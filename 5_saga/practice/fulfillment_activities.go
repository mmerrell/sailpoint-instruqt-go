package fulfillment

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

// FulfillmentActivities groups the forward steps and their compensating
// transactions for the Saga-pattern fulfillment workflow.
type FulfillmentActivities struct{}

// Forward steps — run in order.

// ReserveInventory holds stock for the order. 10% simulated failure rate.
func (a *FulfillmentActivities) ReserveInventory(ctx context.Context, order Order) (string, error) {
	activity.GetLogger(ctx).Info("Reserving inventory", "orderId", order.OrderID)
	if rand.Float64() < 0.1 {
		return "", temporal.NewApplicationError("inventory service unavailable", "InventoryError")
	}
	return fmt.Sprintf("RES-%s-%d", order.OrderID, time.Now().UnixMilli()), nil
}

// ProcessPayment charges the customer. 20% simulated failure rate.
func (a *FulfillmentActivities) ProcessPayment(ctx context.Context, order Order) (string, error) {
	activity.GetLogger(ctx).Info("Processing payment", "orderId", order.OrderID)
	if rand.Float64() < 0.2 {
		return "", temporal.NewApplicationError("payment gateway unavailable", "PaymentError")
	}
	return fmt.Sprintf("PAY-%s-%d", order.OrderID, time.Now().UnixMilli()), nil
}

// DispatchToFulfillment ships the order. 30% simulated failure rate — high on
// purpose, so learners reliably see compensation run within a few tries.
func (a *FulfillmentActivities) DispatchToFulfillment(ctx context.Context, order Order, reservationID string) (string, error) {
	activity.GetLogger(ctx).Info("Dispatching order", "orderId", order.OrderID)
	if rand.Float64() < 0.3 {
		return "", temporal.NewApplicationError("fulfillment API error", "DispatchError")
	}
	return fmt.Sprintf("TRK-%d-%d", len(reservationID), time.Now().UnixMilli()), nil
}

// Compensating transactions — run in reverse order if a later forward step fails.
// There is no compensation for DispatchToFulfillment: if it fails, it never
// completed, so there's nothing to undo there.

// ReleaseInventory undoes ReserveInventory.
func (a *FulfillmentActivities) ReleaseInventory(ctx context.Context, reservationID string) error {
	activity.GetLogger(ctx).Info("Releasing inventory reservation", "reservationId", reservationID)
	// In production this would call the inventory service to free the hold.
	return nil
}

// RefundPayment undoes ProcessPayment.
func (a *FulfillmentActivities) RefundPayment(ctx context.Context, paymentConfirmation string) error {
	activity.GetLogger(ctx).Info("Refunding payment", "paymentConfirmation", paymentConfirmation)
	// In production this would call the payment gateway's refund endpoint.
	return nil
}
