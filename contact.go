package generator

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"image"

	"github.com/jung-kurt/gofpdf"
)

// Contact contact a company informations
type Contact struct {
	Name    string   `json:"name,omitempty" validate:"required,min=1,max=256"`
	Logo    *[]byte  `json:"logo,omitempty"` // Logo byte array
	Address *Address `json:"address,omitempty"`
}

func (c *Contact) appendContactTODoc(x float64, y float64, fill bool, logoAlign string, pdf *gofpdf.Fpdf) float64 {
	pdf.SetXY(x, y)

	// Logo
	if c.Logo != nil {
		// Create filename
		fileName := b64.StdEncoding.EncodeToString([]byte(c.Name))
		// Create reader from logo bytes
		ioReader := bytes.NewReader(*c.Logo)
		// Get image format
		_, format, _ := image.DecodeConfig(bytes.NewReader(*c.Logo))
		// Register image in pdf
		imageInfo := pdf.RegisterImageOptionsReader(fileName, gofpdf.ImageOptions{
			ImageType: format,
		}, ioReader)

		if imageInfo != nil {
			var imageOpt gofpdf.ImageOptions
			imageOpt.ImageType = format

			pdf.ImageOptions(fileName, pdf.GetX(), y, 0, 9, false, imageOpt, 0, "")

			pdf.SetY(y + 20)
		}
	}

	// Name
	if fill {
		pdf.SetFillColor(GreyBgColor[0], GreyBgColor[1], GreyBgColor[2])
	} else {
		pdf.SetFillColor(255, 255, 255)
	}

	// Reset x
	pdf.SetX(x)

	// Name rect
	pdf.Rect(x, pdf.GetY(), 70, 8, "F")

	// Set name
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Cell(40, 8, encodeString(c.Name))
	pdf.SetFont("Helvetica", "", 10)

	if c.Address != nil {
		// Address rect
		var addrRectHeight float64 = 17

		if len(c.Address.Address2) > 0 {
			addrRectHeight = addrRectHeight + 5
		}

		if len(c.Address.Country) == 0 {
			addrRectHeight = addrRectHeight - 5
		}

		pdf.Rect(x, pdf.GetY()+9, 70, addrRectHeight, "F")

		// Set address
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetXY(x, pdf.GetY()+10)
		pdf.MultiCell(70, 5, c.Address.ToString(), "0", "L", false)
	}

	return pdf.GetY()
}

func (c *Contact) appendContactAddressHeaderTODoc(x float64, y float64, pdf *gofpdf.Fpdf) float64 {
	pdf.SetXY(x, y)
	pdf.SetX(x)
	pdf.SetFillColor(255, 255, 255)

	if c.Address != nil {
		// Address rect
		var addrRectHeight float64 = 5
		pdf.Rect(x, pdf.GetY(), 80, addrRectHeight, "F")

		// Set address
		pdf.SetFont("Helvetica", "", 7)
		pdf.SetXY(x, pdf.GetY())
		pdf.MultiCell(70, 5, fmt.Sprintf("%s - %s", encodeString(c.Name), c.Address.ToLineString()), "0", "L", false)
	}

	return pdf.GetY()
}

func (c *Contact) appendCompanyContactToDoc(pdf *gofpdf.Fpdf) float64 {
	x, y, _, _ := pdf.GetMargins()
	c.appendContactAddressHeaderTODoc(x+10, y+26, pdf)
	return c.appendContactTODoc(140, BaseMarginTop, false, "L", pdf)
}

func (c *Contact) appendCustomerContactToDoc(pdf *gofpdf.Fpdf) float64 {
	x, y, _, _ := pdf.GetMargins()
	return c.appendContactTODoc(x+12, y+32, false, "L", pdf)
}
