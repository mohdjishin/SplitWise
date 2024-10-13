package dto

// CreateGroupWithBillResponse represents the response body for creating a group with an associated bill.
// @Description Response model for the creation of a group with an associated bill.
// @Name CreateGroupWithBillResponse
// @Property groupId integer "Id of the created group"
// @Property billId integer "Id of the created bill"
// @Property message string "Success message"
type CreateGroupWithBillResponse struct {
	GroupID uint   `json:"groupId"`
	BillID  uint   `json:"billId"`
	Message string `json:"message"`
}
