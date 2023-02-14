package generator

import "fmt"

// Address represent an address
type Address struct {
	Address    string `json:"address,omitempty" validate:"required"`
	Address2   string `json:"address_2,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Extend     string `json:"extend,omitempty"`
}

// ToString output address as string
// Line break are added for new lines
func (a *Address) ToString() string {
	var addrString string = a.Address

	if len(a.Address2) > 0 {
		addrString += "\n"
		addrString += a.Address2
	}

	if len(a.PostalCode) > 0 {
		addrString += "\n"
		addrString += a.PostalCode
	} else {
		addrString += "\n"
	}

	if len(a.City) > 0 {
		addrString += " "
		addrString += a.City
	}

	if len(a.Country) > 0 {
		addrString += "\n"
		addrString += a.Country
	}

	if len(a.Phone) > 0 {
		addrString += "\n"
		addrString += fmt.Sprintf("Tel.: %s", a.Phone)
	}

	if len(a.Extend) > 0 {
		addrString += "\n"
		addrString += a.Extend
	}

	return encodeString(addrString)
}

func (a *Address) ToLineString() string {
	var addrString string = a.Address

	if len(a.Address2) > 0 {
		addrString += " - "
		addrString += a.Address2
	}

	if len(a.PostalCode) > 0 {
		addrString += " - "
		addrString += a.PostalCode
	}

	if len(a.City) > 0 {
		addrString += " "
		addrString += a.City
	}

	return encodeString(addrString)
}
