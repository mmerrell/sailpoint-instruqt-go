package fulfillment

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "fulfillment-tasks"

// FulfillmentWorkflow orchestrates order fulfillment using the Saga pattern:
// reserve inventory, charge payment, dispatch. If a step fails partway
// through, a deferred compensation block undoes whichever steps actually
// completed, in reverse order.
func FulfillmentWorkflow(ctx workflow.Context, order Order) (result OrderResult, err error) {
	log := workflow.GetLogger(ctx)
	log.Info("Processing order", "orderId", order.OrderID)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	fa := &FulfillmentActivities{}

	// State tracked across the deferred compensation block below — each stays
	// "" until its corresponding forward step actually completes.
	var reservationID, paymentConfirmation string

	// Go has no try/catch. This deferred func plays the same role as a catch
	// block: it always runs on the way out of the function, and only takes
	// action if the workflow is returning with a non-nil error. Because it's
	// a closure over reservationID/paymentConfirmation, it sees exactly which
	// steps completed no matter which one failed.
	defer func() {
		if err == nil {
			return
		}
		log.Warn("Order failed — running compensations", "orderId", order.OrderID, "cause", err)

		if paymentConfirmation != "" {
			if cErr := workflow.ExecuteActivity(ctx, fa.RefundPayment, paymentConfirmation).Get(ctx, nil); cErr != nil {
				log.Error("Compensation failed: RefundPayment", "error", cErr)
			} else {
				log.Info("Payment refunded", "paymentConfirmation", paymentConfirmation)
			}
		}
		if reservationID != "" {
			if cErr := workflow.ExecuteActivity(ctx, fa.ReleaseInventory, reservationID).Get(ctx, nil); cErr != nil {
				log.Error("Compensation failed: ReleaseInventory", "error", cErr)
			} else {
				log.Info("Inventory released", "reservationId", reservationID)
			}
		}

		// The saga has already handled the failure — return a FAILED result,
		// not a workflow error.
		result = OrderResult{
			OrderID:             order.OrderID,
			Status:              "FAILED",
			ReservationID:       reservationID,
			PaymentConfirmation: paymentConfirmation,
		}
		err = nil
	}()

	if err = workflow.ExecuteActivity(ctx, fa.ReserveInventory, order).Get(ctx, &reservationID); err != nil {
		return
	}
	log.Info("Inventory reserved", "reservationId", reservationID)

	if err = workflow.ExecuteActivity(ctx, fa.ProcessPayment, order).Get(ctx, &paymentConfirmation); err != nil {
		return
	}
	log.Info("Payment confirmed", "paymentConfirmation", paymentConfirmation)

	var trackingNumber string
	if err = workflow.ExecuteActivity(ctx, fa.DispatchToFulfillment, order, reservationID).Get(ctx, &trackingNumber); err != nil {
		return
	}
	log.Info("Order dispatched", "trackingNumber", trackingNumber)

	result = OrderResult{
		OrderID:             order.OrderID,
		Status:              "FULFILLED",
		ReservationID:       reservationID,
		PaymentConfirmation: paymentConfirmation,
		TrackingNumber:      trackingNumber,
	}
	return
}
