package generator

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

// Build pdf document from data provided
func (d *Document) Build() (*gofpdf.Fpdf, error) {
	// Validate document data
	err := d.Validate()
	if err != nil {
		return nil, err
	}

	// Build base doc
	d.pdf.SetMargins(BaseMargin, BaseMarginTop, BaseMargin)
	d.pdf.SetXY(10, 10)
	d.pdf.SetTextColor(BaseTextColor[0], BaseTextColor[1], BaseTextColor[2])

	// Set header
	if d.Header != nil {
		err = d.Header.applyHeader(d, d.pdf)

		if err != nil {
			return nil, err
		}
	}

	// Set footer
	if d.Footer != nil {
		err = d.Footer.applyFooter(d, d.pdf)

		if err != nil {
			return nil, err
		}
	}

	// Add first page
	d.pdf.AddPage()

	// Load font
	d.pdf.SetFont("Helvetica", "", 12)

	// Append document title
	d.appendTitle(d.pdf)

	// Append document metas (ref & version)
	d.appendMetas(d.pdf)

	// Append company contact to doc
	companyBottom := d.Company.appendCompanyContactToDoc(d.pdf)

	// Append customer contact to doc
	customerBottom := d.Customer.appendCustomerContactToDoc(d.pdf)

	if customerBottom > companyBottom {
		d.pdf.SetXY(10, customerBottom)
	} else {
		d.pdf.SetXY(10, companyBottom)
	}

	// Append description
	d.appendDescription(d.pdf)

	// Append items
	d.appendItems(d.pdf)

	// Check page height (total bloc height = 30, 45 when doc discount)
	offset := d.pdf.GetY() + 35
	if d.Discount != nil {
		offset += 15
	}
	if offset > MaxPageHeight {
		d.pdf.AddPage()
	}

	// Append notes
	d.appendNotes(d.pdf)

	// Append total
	d.appendTotal(d.pdf)

	// Append payment term
	d.appendPaymentTerm(d.pdf)

	// Append js to autoprint if AutoPrint == true
	if d.Options.AutoPrint {
		d.pdf.SetJavascript("print(true);")
	}

	return d.pdf, nil
}

func (d *Document) appendTitle(pdf *gofpdf.Fpdf) {
	title := d.typeAsString()

	x, _, _, _ := pdf.GetMargins()

	// Set x y
	pdf.SetXY(x, BaseMarginTop+68)

	// Draw rect
	//pdf.SetFillColor(DarkBgColor[0], DarkBgColor[1], DarkBgColor[2])
	//pdf.Rect(x, BaseMarginTop+60, 80, 10, "F")

	// Draw text
	pdf.SetFont("Helvetica", "", 14)
	pdf.CellFormat(x+15, 10, encodeString(title), "0", 0, "C", false, 0, "")
}

func (d *Document) appendMetas(pdf *gofpdf.Fpdf) {
	// Append ref
	refString := fmt.Sprintf("%s: %s", d.Options.TextRefTitle, d.Ref)

	offset := 68.0
	pdf.SetXY(140, BaseMarginTop+offset)
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(50, 4, encodeString(refString), "0", 0, "L", false, 0, "")

	// Append version
	if len(d.Version) > 0 {
		offset = offset + 4
		versionString := fmt.Sprintf("%s: %s", d.Options.TextVersionTitle, d.Version)
		pdf.SetXY(140, BaseMarginTop+offset)
		pdf.SetFont("Helvetica", "", 8)
		pdf.CellFormat(50, 4, encodeString(versionString), "0", 0, "R", false, 0, "")
	}

	// Append date
	date := time.Now().Format("01.02.2006")
	if len(d.Date) > 0 {
		date = d.Date
	}
	offset = offset + 4
	dateString := fmt.Sprintf("%s: %s", d.Options.TextDateTitle, date)
	pdf.SetXY(140, BaseMarginTop+offset)
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(50, 4, encodeString(dateString), "0", 0, "L", false, 0, "")
}

func (d *Document) appendDescription(pdf *gofpdf.Fpdf) {
	if len(d.Description) > 0 {
		pdf.SetY(pdf.GetY() + 30)
		pdf.SetFont("Helvetica", "", 10)
		pdf.MultiCell(190, 5, encodeString(d.Description), "B", "L", false)
	}
}

func (d *Document) drawsTableTitles(pdf *gofpdf.Fpdf) {
	// Draw table titles
	pdf.SetX(10)
	pdf.SetY(pdf.GetY() + 5)
	pdf.SetFont("Helvetica", "B", 8)

	// Draw rec
	pdf.SetFillColor(GreyBgColor[0], GreyBgColor[1], GreyBgColor[2])
	pdf.Rect(10, pdf.GetY(), 190, 6, "F")

	// Name
	pdf.SetX(ItemColNameOffset)
	pdf.CellFormat(
		ItemColUnitPriceOffset-ItemColNameOffset,
		6,
		encodeString(d.Options.TextItemsNameTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// Unit price
	if len(d.Options.TextItemsUnitCostTitle) > 1 {
		pdf.SetX(ItemColUnitPriceOffset)
		pdf.CellFormat(
			ItemColQuantityOffset-ItemColUnitPriceOffset,
			6,
			encodeString(d.Options.TextItemsUnitCostTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// Quantity
	if len(d.Options.TextItemsQuantityTitle) > 1 {
		pdf.SetX(ItemColQuantityOffset)
		pdf.CellFormat(
			ItemColTaxOffset-ItemColQuantityOffset,
			6,
			encodeString(d.Options.TextItemsQuantityTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// Total HT
	if len(d.Options.TextItemsTotalHTTitle) > 1 {
		pdf.SetX(ItemColTotalHTOffset)
		pdf.CellFormat(
			ItemColTaxOffset-ItemColTotalHTOffset,
			6,
			encodeString(d.Options.TextItemsTotalHTTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// Tax
	if len(d.Options.TextItemsTaxTitle) > 1 {
		pdf.SetX(ItemColTaxOffset)
		pdf.CellFormat(
			ItemColProductNumberOffset-ItemColTaxOffset,
			6,
			encodeString(d.Options.TextItemsTaxTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// ProductNumber
	if len(d.Options.TextItemsProductNumberTitle) > 1 {
		pdf.SetX(ItemColProductNumberOffset)
		pdf.CellFormat(
			ItemColTotalTTCOffset-ItemColProductNumberOffset,
			6,
			encodeString(d.Options.TextItemsProductNumberTitle),
			"0",
			0,
			"",
			false,
			0,
			"",
		)
	}

	// TOTAL TTC
	if len(d.Options.TextItemsTotalTTCTitle) > 1 {
		pdf.SetX(ItemColTotalTTCOffset)
		pdf.CellFormat(190-ItemColTotalTTCOffset, 6, encodeString(d.Options.TextItemsTotalTTCTitle), "0", 0, "", false, 0, "")
	}
}

func (d *Document) appendItems(pdf *gofpdf.Fpdf) {
	d.drawsTableTitles(pdf)

	pdf.SetX(10)
	pdf.SetY(pdf.GetY() + 8)
	pdf.SetFont("Helvetica", "", 8)

	for i := 0; i < len(d.Items); i++ {
		item := d.Items[i]

		// Check item tax
		if item.Tax == nil {
			item.Tax = d.DefaultTax
		}

		// Append to pdf
		item.appendColTo(d.Options, pdf)

		if pdf.GetY() > MaxPageHeight {
			// Add page
			pdf.AddPage()
			d.drawsTableTitles(pdf)
			pdf.SetFont("Helvetica", "", 8)
		}

		pdf.SetX(10)
		pdf.SetY(pdf.GetY() + 6)
	}
}

func (d *Document) appendNotes(pdf *gofpdf.Fpdf) {
	if len(d.Notes) == 0 {
		return
	}

	currentY := pdf.GetY()

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetX(BaseMargin)
	pdf.SetRightMargin(100)
	pdf.SetY(currentY + 10)

	_, lineHt := pdf.GetFontSize()
	html := pdf.HTMLBasicNew()
	html.Write(lineHt, encodeString(d.Notes))

	pdf.SetRightMargin(BaseMargin)
	pdf.SetY(currentY)
}

func (d *Document) appendTotal(pdf *gofpdf.Fpdf) {
	ac := accounting.Accounting{
		Symbol:    encodeString(d.Options.CurrencySymbol),
		Precision: d.Options.CurrencyPrecision,
		Thousand:  d.Options.CurrencyThousand,
		Decimal:   d.Options.CurrencyDecimal,
	}

	// Get total (without tax)
	total, _ := decimal.NewFromString("0")

	for _, item := range d.Items {
		total = total.Add(item.totalWithoutTaxAndWithDiscount())
	}

	// Apply document discount
	totalWithDiscount := decimal.NewFromFloat(0)
	if d.Discount != nil {
		discountType, discountNumber := d.Discount.getDiscount()

		if discountType == "amount" {
			totalWithDiscount = total.Sub(discountNumber)
		} else {
			// Percent
			toSub := total.Mul(discountNumber.Div(decimal.NewFromFloat(100)))
			totalWithDiscount = total.Sub(toSub)
		}
	}

	// Tax
	totalTax := decimal.NewFromFloat(0)
	if d.Discount == nil {
		for _, item := range d.Items {
			totalTax = totalTax.Add(item.taxWithDiscount())
		}
	} else {
		discountType, discountAmount := d.Discount.getDiscount()
		discountPercent := discountAmount
		if discountType == "amount" {
			// Get percent from total discounted
			discountPercent = discountAmount.Mul(decimal.NewFromFloat(100)).Div(totalWithDiscount)
		}

		for _, item := range d.Items {
			if item.Tax != nil {
				taxType, taxAmount := item.Tax.getTax()
				if taxType == "amount" {
					// If tax type is amount, juste add amount to tax
					totalTax = totalTax.Add(taxAmount)
				} else {
					// Else, remove doc discount % from item total without tax and item discount
					itemTotal := item.totalWithoutTaxAndWithDiscount()
					toSub := discountPercent.Mul(itemTotal).Div(decimal.NewFromFloat(100))
					itemTotalDiscounted := itemTotal.Sub(toSub)

					// Then recompute tax on itemTotalDiscounted
					itemTaxDiscounted := taxAmount.Mul(itemTotalDiscounted).Div(decimal.NewFromFloat(100))

					totalTax = totalTax.Add(itemTaxDiscounted)
				}
			}
		}
	}

	// finalTotal
	totalWithTax := total.Add(totalTax)
	if d.Discount != nil {
		totalWithTax = totalWithDiscount.Add(totalTax)
	}

	pdf.SetY(pdf.GetY() + 10)
	pdf.SetFont("Helvetica", "", LargeTextFontSize)
	pdf.SetTextColor(BaseTextColor[0], BaseTextColor[1], BaseTextColor[2])

	if d.Options.TextTotalTotal != " " {
		// Draw TOTAL HT title
		pdf.SetX(120)
		pdf.SetFillColor(DarkBgColor[0], DarkBgColor[1], DarkBgColor[2])
		pdf.Rect(120, pdf.GetY(), 40, 10, "F")
		pdf.CellFormat(38, 10, encodeString(d.Options.TextTotalTotal), "0", 0, "R", false, 0, "")

		// Draw TOTAL HT amount
		pdf.SetX(162)
		pdf.SetFillColor(GreyBgColor[0], GreyBgColor[1], GreyBgColor[2])
		pdf.Rect(160, pdf.GetY(), 40, 10, "F")
		pdf.CellFormat(40, 10, ac.FormatMoneyDecimal(total), "0", 0, "L", false, 0, "")
	}

	if d.Discount != nil && len(d.Options.TextTotalDiscounted) > 1 {
		baseY := pdf.GetY() + 10

		// Draw DISCOUNTED title
		pdf.SetXY(120, baseY)
		pdf.SetFillColor(DarkBgColor[0], DarkBgColor[1], DarkBgColor[2])
		pdf.Rect(120, pdf.GetY(), 40, 15, "F")

		// title
		pdf.CellFormat(38, 7.5, encodeString(d.Options.TextTotalDiscounted), "0", 0, "BR", false, 0, "")

		// description
		pdf.SetXY(120, baseY+7.5)
		pdf.SetFont("Helvetica", "", BaseTextFontSize)
		pdf.SetTextColor(GreyTextColor[0], GreyTextColor[1], GreyTextColor[2])

		var descString bytes.Buffer
		discountType, discountAmount := d.Discount.getDiscount()
		if discountType == "percent" {
			descString.WriteString("-")
			descString.WriteString(discountAmount.String())
			descString.WriteString(" % / -")
			descString.WriteString(ac.FormatMoneyDecimal(total.Sub(totalWithDiscount)))
		} else {
			descString.WriteString("-")
			descString.WriteString(ac.FormatMoneyDecimal(discountAmount))
			descString.WriteString(" / -")
			descString.WriteString(discountAmount.Mul(decimal.NewFromFloat(100)).Div(total).StringFixed(2))
			descString.WriteString(" %")
		}

		pdf.CellFormat(38, 7.5, descString.String(), "0", 0, "TR", false, 0, "")

		pdf.SetFont("Helvetica", "", LargeTextFontSize)
		pdf.SetTextColor(BaseTextColor[0], BaseTextColor[1], BaseTextColor[2])

		// Draw DISCOUNT amount
		pdf.SetY(baseY)
		pdf.SetX(162)
		pdf.SetFillColor(GreyBgColor[0], GreyBgColor[1], GreyBgColor[2])
		pdf.Rect(160, pdf.GetY(), 40, 15, "F")
		pdf.CellFormat(40, 15, ac.FormatMoneyDecimal(totalWithDiscount), "0", 0, "L", false, 0, "")
		pdf.SetY(pdf.GetY() + 15)
	} else {
		pdf.SetY(pdf.GetY() + 10)
	}

	// Draw TAX title
	if len(d.Options.TextTotalTax) > 1 {
		pdf.SetX(120)
		pdf.SetFillColor(DarkBgColor[0], DarkBgColor[1], DarkBgColor[2])
		pdf.Rect(120, pdf.GetY(), 40, 10, "F")
		pdf.CellFormat(38, 10, encodeString(d.Options.TextTotalTax), "0", 0, "R", false, 0, "")

		// Draw TAX amount
		pdf.SetX(162)
		pdf.SetFillColor(GreyBgColor[0], GreyBgColor[1], GreyBgColor[2])
		pdf.Rect(160, pdf.GetY(), 40, 10, "F")
		pdf.CellFormat(40, 10, ac.FormatMoneyDecimal(totalTax), "0", 0, "L", false, 0, "")
	}

	// Draw TOTAL TTC title
	if len(d.Options.TextTotalWithTax) > 1 {
		pdf.SetY(pdf.GetY() + 10)
		pdf.SetX(120)
		pdf.SetFillColor(DarkBgColor[0], DarkBgColor[1], DarkBgColor[2])
		pdf.Rect(120, pdf.GetY(), 40, 10, "F")
		pdf.CellFormat(38, 10, encodeString(d.Options.TextTotalWithTax), "0", 0, "R", false, 0, "")

		// Draw TOTAL TTC amount
		pdf.SetX(162)
		pdf.SetFillColor(GreyBgColor[0], GreyBgColor[1], GreyBgColor[2])
		pdf.Rect(160, pdf.GetY(), 40, 10, "F")
		pdf.CellFormat(40, 10, ac.FormatMoneyDecimal(totalWithTax), "0", 0, "L", false, 0, "")
	}
}

func (d *Document) appendPaymentTerm(pdf *gofpdf.Fpdf) {
	if len(d.BankingDetails) > 0 {
		bankingDetailsString := fmt.Sprintf("%s", encodeString(d.BankingDetails))
		pdf.SetY(pdf.GetY() + 15)

		pdf.SetX(120)
		pdf.SetFont("Helvetica", "", 10)
		pdf.CellFormat(80, 4, bankingDetailsString, "0", 0, "R", false, 0, "")
	}
	if len(d.PaymentTerm) > 0 {
		paymentTermString := fmt.Sprintf("%s: %s", encodeString(d.Options.TextPaymentTermTitle), encodeString(d.PaymentTerm))
		pdf.SetY(pdf.GetY() + 4)

		pdf.SetX(120)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(80, 4, paymentTermString, "0", 0, "R", false, 0, "")
	}
}
