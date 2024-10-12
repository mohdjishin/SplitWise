package dto

// GetGroupReportRequest represents the request body for retrieving a group report.
// @Description Request model for generating a group report based on date range.
// @Name GetGroupReportRequest
// @Property from string true "Start date in the format YYYY-MM-DD" example("2024-10-01")
// @Property to
type GetGroupReportRequest struct {
	From *string `json:"from" validate:"omitempty,dateFormat"`
	To   *string `json:"to"   validate:"omitempty,dateFormat"`
}
