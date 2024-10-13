package dto

// MarkPaymentRequest represents the request body for marking a payment.
// @Description Mark a payment for a specific group
// @Param request body MarkPaymentRequest true "Mark Payment Request"
// &Param request body MarkPaymentRequest true "Mark Payment Request"
// @Example { "groupId": 1, "remarks": "Paid for dinner" }
type MarkPaymentRequest struct {
	GroupID uint   `json:"groupId" validate:"required"`
	Remarks string `json:"remarks"` // Optional* remarks for the payment
}

// MarkPaymentResponse represents the response returned after marking a payment.
// @Description Response for marking a payment
// @Success 200 {object} MarkPaymentResponse "Payment marked successfully"
// @Example { "message": "Payment marked successfully" }
type MarkPaymentResponse struct {
	Message string `json:"message"`
}
