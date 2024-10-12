package pdf

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/mohdjishin/SplitWise/internal/models/dto"
)

func GeneratePDFReport(groupsInfo []dto.Group, startDate, endDate time.Time, ownerName string) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Tabloid", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, fmt.Sprintf("Monthly/Weekly Report for Groups Owned by %s", ownerName))
	pdf.Ln(10)

	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 10, fmt.Sprintf("Report Period: %s - %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 7)
	headers := []string{"Group ID", "Group Name", "Total Amount", "Share Per Member", "Paid Amount", "Members", "Status", "Last Bill Date", "Notes"}
	headerWidths := []float64{15, 35, 25, 25, 15, 20, 25, 25, 30}

	for i, header := range headers {
		pdf.CellFormat(headerWidths[i], 8, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 7)

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
			"N/A", //TODO: have to add later
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
				pdf.MultiCell(headerWidths[i], 5, data, "1", "C", false)
				pdf.SetXY(x+headerWidths[i], y)
			} else {
				pdf.CellFormat(headerWidths[i], rowHeight, data, "1", 0, "C", false, 0, "")
			}
		}
		pdf.Ln(rowHeight)
	}

	return pdf
}
