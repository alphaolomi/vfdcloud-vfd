package vfd

import (
	"fmt"
	"math"
)

const (
	StandardVATID        = "A"
	StandardVATRATE      = 18.00
	StandardVATCODE      = 1
	SpecialVATID         = "B"
	SpecialVATRATE       = 10.00
	SpecialVATCODE       = 2
	ZeroVATID            = "C"
	ZeroVATRATE          = 0.00
	ZeroVATCODE          = 3
	SpecialReliefVATID   = "D"
	SpecialReliefVATRATE = 0.00
	SpecialReliefVATCODE = 4
	ExemptedVATID        = "E"
	ExemptedVATRATE      = 0.00
	ExemptedVATCODE      = 5
)

type (
	vat struct {
		ID         string // ID is a character that identifies the vat it can be A,B,C,D or E
		Code       int64  // Code is a number that identifies the vat it can be 0,1,2,3 or 4
		Name       string
		Percentage float64
	}
)

var (
	standardVAT = vat{
		ID:         StandardVATID,
		Code:       StandardVATCODE,
		Name:       "Standard vat",
		Percentage: StandardVATRATE,
	}
	specialVAT = vat{
		ID:         SpecialVATID,
		Code:       SpecialVATCODE,
		Name:       "Special vat",
		Percentage: SpecialVATRATE,
	}
	zeroVAT = vat{
		ID:         ZeroVATID,
		Code:       ZeroVATCODE,
		Name:       "Zero vat",
		Percentage: ZeroVATRATE,
	}
	specialReliefVAT = vat{
		ID:         SpecialReliefVATID,
		Code:       SpecialReliefVATCODE,
		Name:       "Special Relief vat",
		Percentage: SpecialReliefVATRATE,
	}
	exemptedVAT = vat{
		ID:         ExemptedVATID,
		Code:       ExemptedVATCODE,
		Name:       "Exempted vat",
		Percentage: ExemptedVATRATE,
	}
)

// amount returns the amount deducted from the price of the product
// of a certain vat category. Answer is rounded to 2 decimal places.
func (v *vat) amount(price float64) float64 {
	return math.Floor((price*v.Percentage/100)*100) / 100
}

func parseTaxCode(code int64) vat {
	switch code {
	case 1:
		return standardVAT
	case 2:
		return specialVAT
	case 3:
		return zeroVAT
	case 4:
		return specialReliefVAT
	case 5:
		return exemptedVAT
	default:
		return standardVAT
	}
}

// ValueAddedTaxRate returns the vat rate of a certain vat category
func ValueAddedTaxRate(taxCode int64) float64 {
	vat := parseTaxCode(taxCode)
	return vat.Percentage
}

// ValueAddedTaxID returns the vat id of a certain vat category
// It returns "A" for standard vat, "B" for special vat,
// "C" for zero vat,"D" for special relief and "E" for
// exempted vat.
func ValueAddedTaxID(taxCode int64) string {
	vat := parseTaxCode(taxCode)
	return vat.ID
}

func TaxAmount(taxCode int64, price float64) float64 {
	vat := parseTaxCode(taxCode)
	return vat.amount(price)
}

// ReportTaxRateID creates a string that contains the vat rate and the vat id
// of a certain vat category. It returns "A-18.00" for standard vat,
// "B-10.00" for special vat, "C-0.00" for zero vat and so on. The ID is then
// used in Z Report to indicate the vat rate and the vat id.
func ReportTaxRateID(taxCode int64) string {
	vat := parseTaxCode(taxCode)
	return fmt.Sprintf("%s-%.2f", vat.ID, vat.Percentage)
}
