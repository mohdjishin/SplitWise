package dto

// PendingPaymentsWithTotalResponse  represents the response for listing pending payments.
// @Description Response model for listing pending payments with total amount.
// @Name PendingPaymentsWithTotalResponse
// @Property pendingPayments []PendingPayments true "List of pending payments"
// @Property totalAmount float64 true "Total amount of pending payments"
type PendingPaymentsWithTotalResponse struct {
	PendingPayments []PendingPayments `json:"pendingPayments"`
	TotalAmount     float64           `json:"totalAmount"`
}

type PendingPayments struct {
	GroupID   uint    `json:"groupId"`
	GroupName string  `json:"groupName"`
	BillID    uint    `json:"billId"`
	Amount    float64 `json:"amount"`
}
