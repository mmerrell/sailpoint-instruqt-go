package fulfillment

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "fulfillment-tasks"

// FulfillmentWorkflow orchestrates order fulfillment using the Saga pattern:
// reserve inventory, charge payment, dispatch. If a step fails partway
// through, compensations undo whichever steps actually completed, in
// reverse order.
func FulfillmentWorkflow(ctx workflow.Context, order Order) (result OrderResult, err error) {
	log := workflow.GetLogger(ctx)
	log.Info("Processing order", "orderId", order.OrderID)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	fa := &FulfillmentActivities{}

	// State tracked across the compensation block — each stays "" until its
	// corresponding forward step actually completes. Declared here (not
	// inside an if-block) so both the forward steps and the compensation
	// logic below can see them.
	var reservationID, paymentConfirmation string

	// TODO Part C: Go has no try/catch. Add a deferred func here — right
	// after the state above is declared — that plays the same role as
	// Java's catch block: it always runs on the way out of the function,
	// and only takes action if the workflow is returning with a non-nil
	// `err`. Because it's a closure over reservationID/paymentConfirmation,
	// it will see exactly which steps completed no matter which one failed.
	//
	// Guard each compensation on whether the corresponding step actually
	// completed:
	//
	//   defer func() {
	//       if err == nil {
	//           return
	//       }
	//       if paymentConfirmation != "" {
	//           _ = workflow.ExecuteActivity(ctx, fa.RefundPayment, paymentConfirmation).Get(ctx, nil)
	//       }
	//       if reservationID != "" {
	//           _ = workflow.ExecuteActivity(ctx, fa.ReleaseInventory, reservationID).Get(ctx, nil)
	//       }
	//       result = OrderResult{
	//           OrderID:             order.OrderID,
	//           Status:              "FAILED",
	//           ReservationID:       reservationID,
	//           PaymentConfirmation: paymentConfirmation,
	//       }
	//       err = nil // the saga already handled the failure — return a result, not an error
	//   }()

	// TODO Part B: Call the three forward steps in order, assigning into the
	// named return `err` and returning immediately if any of them fails:
	//
	//   if err = workflow.ExecuteActivity(ctx, fa.ReserveInventory, order).Get(ctx, &reservationID); err != nil {
	//       return
	//   }
	//   if err = workflow.ExecuteActivity(ctx, fa.ProcessPayment, order).Get(ctx, &paymentConfirmation); err != nil {
	//       return
	//   }
	//   var trackingNumber string
	//   if err = workflow.ExecuteActivity(ctx, fa.DispatchToFulfillment, order, reservationID).Get(ctx, &trackingNumber); err != nil {
	//       return
	//   }
	//
	// On success, set `result` to a "FULFILLED" OrderResult (include
	// trackingNumber) and return.

	_ = fa                    // remove once you call fa.ReserveInventory etc. above
	_ = reservationID         // remove once Part B assigns into it
	_ = paymentConfirmation   // remove once Part B assigns into it

	// Placeholder — remove once Part B and Part C are implemented.
	result = OrderResult{OrderID: order.OrderID, Status: "NOT_IMPLEMENTED"}
	return
}
