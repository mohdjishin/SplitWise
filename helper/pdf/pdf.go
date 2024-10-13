package pdf

import (
	"fmt"

	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/mohdjishin/SplitWise/internal/models/dto"
)

func GenerateOwnerGroupsReportPDF(groupsInfo []dto.Group, startDate, endDate time.Time, ownerName string) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 51, 102)
	pdf.Cell(0, 10, fmt.Sprintf("Monthly/Weekly Report for Groups Owned by %s", ownerName))
	pdf.Ln(10)

	pdf.SetFont("Arial", "I", 9)
	pdf.SetTextColor(33, 33, 33)
	pdf.Cell(0, 10, fmt.Sprintf("Report Period: %s - %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	pdf.Ln(12)

	pdf.SetFillColor(200, 200, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 7)
	headers := []string{"Group ID", "Group Name", "Total Amount", "Share Per Member", "Paid Amount", "Members", "Status", "Last Bill Date"}
	headerWidths := []float64{15, 35, 25, 25, 15, 20, 25, 25}

	for i, header := range headers {
		pdf.CellFormat(headerWidths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 7)
	pdf.SetFillColor(255, 255, 255)

	for _, group := range groupsInfo {
		lastBillDate := "N/A"
		if !group.Bills.Date.IsZero() {
			lastBillDate = group.Bills.Date.Format("2006-01-02")
		}

		rowData := []string{
			group.ID,
			group.Name,
			fmt.Sprintf("$%.2f", group.TotalAmount),
			fmt.Sprintf("$%.2f", group.PerUserSplitAmount),
			fmt.Sprintf("$%.2f", group.PaidAmount),
			fmt.Sprintf("%d", group.Members),
			group.Status,
			lastBillDate,
		}

		maxLines := 1
		for i, data := range rowData {
			numLines := pdf.GetStringWidth(data) / headerWidths[i]
			if numLines > float64(maxLines) {
				maxLines = int(numLines)
			}
		}
		rowHeight := 5.0 * float64(maxLines)

		for i, data := range rowData {
			if i == 1 || i == 8 {
				x, y := pdf.GetX(), pdf.GetY()
				pdf.MultiCell(headerWidths[i], 5, data, "1", "C", true)
				pdf.SetXY(x+headerWidths[i], y)
			} else {
				pdf.CellFormat(headerWidths[i], rowHeight, data, "1", 0, "C", true, 0, "")
			}
		}
		pdf.Ln(rowHeight)
	}

	return pdf
}

func GenerateGroupDetailedReportPDF(report dto.GroupReportRequest) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.SetTextColor(0, 51, 102)
	pdf.Cell(0, 10, fmt.Sprintf("Report for Group: %s", report.Group.Name))
	pdf.Ln(12)

	pdf.SetFillColor(200, 200, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "Owner and Bill Details", "", 1, "L", true, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(33, 33, 33)
	pdf.Cell(0, 10, fmt.Sprintf("Owner: %s", report.UserInfo[report.Group.CreatedBy]))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Bill Name: %s", report.Bill.Name))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Total Amount: $%.2f", report.Bill.Amount))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Group ID: %d", report.Bill.GroupID))
	pdf.Ln(12)

	pdf.SetFillColor(200, 200, 255)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "Group Details", "", 1, "L", true, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Group ID: %d", report.Group.ID))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Created At: %s", report.Group.CreatedAt.Format("2006-01-02 15:04:05")))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Updated At: %s", report.Group.UpdatedAt.Format("2006-01-02 15:04:05")))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Total Amount: $%.2f", report.Group.TotalAmount))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Paid Amount: $%.2f", report.Group.PaidAmount))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Status: %s", report.Group.Status))
	pdf.Ln(12)

	pdf.SetFillColor(200, 200, 255)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "Payment History", "", 1, "L", true, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 10)
	for _, history := range report.History {
		pdf.Cell(0, 10, fmt.Sprintf("Paid By: %s, Amount: $%.2f, Paid At: %s",
			history.PaidBy, history.Amount, history.PaidAt.Format("2006-01-02 15:04:05")))
		pdf.Ln(6)
	}
	pdf.Ln(6)
	if len(report.Members) > 0 {
		pdf.SetFillColor(200, 200, 255)
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 10, "Members", "", 1, "L", true, 0, "")
		pdf.Ln(2)
	}

	pdf.SetFont("Arial", "", 10)
	for _, member := range report.Members {
		pdf.Cell(0, 10, fmt.Sprintf("Member Name: %s, Split Amount: $%.2f, Has Paid: %t",
			report.UserInfo[member.UserID], member.SplitAmount, member.HasPaid))
		pdf.Ln(6)
	}

	return pdf
}
