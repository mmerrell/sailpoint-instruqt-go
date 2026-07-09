package fulfillment

// Order is the input to the fulfillment workflow.
type Order struct {
	OrderID     string  `json:"orderId"`
	CustomerID  string  `json:"customerId"`
	ItemSKU     string  `json:"itemSku"`
	Quantity    int     `json:"quantity"`
	TotalAmount float64 `json:"totalAmount"`
}

// OrderResult is returned when the order is fully processed — either
// FULFILLED or FAILED (with completed compensations already applied).
type OrderResult struct {
	OrderID             string `json:"orderId"`
	Status              string `json:"status"`
	ReservationID       string `json:"reservationId"`
	PaymentConfirmation string `json:"paymentConfirmation"`
	TrackingNumber      string `json:"trackingNumber"`
}
