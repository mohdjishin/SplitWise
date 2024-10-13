package dto

// PendingPaymentsWithTotalResponse  represents the response for listing pending payments.
// @Description Response model for listing pending payments with total amount.
// @Name PendingPaymentsWithTotalResponse
// @Property pendingPayments []PendingPayments true "List of pending payments"
// @Property totalAmount float64 true "Total amount of pending payments"
// @Property message string false "any message"
type PendingPaymentsWithTotalResponse struct {
	PendingPayments []PendingPayments `json:"pendingPayments,omitempty"`
	TotalAmount     float64           `json:"totalAmount,omitempty"`
	Message         string            `json:"message,omitempty"`
}

type PendingPayments struct {
	GroupID   uint    `json:"groupId"`
	GroupName string  `json:"groupName"`
	BillID    uint    `json:"billId"`
	Amount    float64 `json:"amount"`
}
