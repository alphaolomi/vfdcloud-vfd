package vfd

import (
	"fmt"
	"math"
)

// 1=Standard (A-18)
// 2=Special rate(B-10)
// 3=Zero Rated (C-0)
// 4=Special Relief(D-SR)
// 5=Exempted(E-EX)
const (
	StandardVATID        = "A"
	StandardVATRATE      = 18.00
	StandardVATCode      = 1
	SpecialVATID         = "B"
	SpecialVATRATE       = 10.00
	SpecialVATCode       = 2
	ZeroVATID            = "C"
	ZeroVATRATE          = 0.00
	ZeroVATCode          = 3
	SpecialReliefVATID   = "D"
	SpecialReliefVATRATE = 0.00
	SpecialReliefVATCode = 4
	ExemptedVATID        = "E"
	ExemptedVATRATE      = 0.00
	ExemptedVATCode      = 5
)

type (
	VAT struct {
		ID         string // ID is a character that identifies the VAT it can be A,B,C,D or E
		Code       int64  // Code is a number that identifies the VAT it can be 0,1,2,3 or 4
		Name       string
		Percentage float64
	}
)

var (
	standardVAT = VAT{
		ID:         StandardVATID,
		Code:       StandardVATCode,
		Name:       "Standard VAT",
		Percentage: StandardVATRATE,
	}
	specialVAT = VAT{
		ID:         SpecialVATID,
		Code:       SpecialVATCode,
		Name:       "Special VAT",
		Percentage: SpecialVATRATE,
	}
	zeroVAT = VAT{
		ID:         ZeroVATID,
		Code:       ZeroVATCode,
		Name:       "Zero VAT",
		Percentage: ZeroVATRATE,
	}
	specialReliefVAT = VAT{
		ID:         SpecialReliefVATID,
		Code:       SpecialReliefVATCode,
		Name:       "Special Relief VAT",
		Percentage: SpecialReliefVATRATE,
	}
	exemptedVAT = VAT{
		ID:         ExemptedVATID,
		Code:       ExemptedVATCode,
		Name:       "Exempted VAT",
		Percentage: ExemptedVATRATE,
	}
)

// Amount returns the amount deducted from the price of the product
// of a certain VAT category. Answer is rounded to 2 decimal places.
func (v *VAT) Amount(price float64) float64 {
	return math.Floor((price*v.Percentage/100)*100) / 100
}

func ParseVATCode(code int64) VAT {
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

func ParseVATID(id string) VAT {
	switch id {
	case "A":
		return standardVAT
	case "B":
		return specialVAT
	case "C":
		return zeroVAT
	case "D":
		return specialReliefVAT
	case "E":
		return exemptedVAT
	default:
		return standardVAT
	}
}

// ReportVATRATE ...
func ReportVATRATEFromCode(vatCode int64) string {
	vat := ParseVATCode(vatCode)
	return fmt.Sprintf("%s-%.2f", vat.ID, vat.Percentage)
}

// ReportVATRATE ...
func ReportVATRATEFromID(vatID string) string {
	vat := ParseVATID(vatID)
	return fmt.Sprintf("%s-%.2f", vat.ID, vat.Percentage)
}

func TaxAmount(vatCode int64, price float64) float64 {
	vat := ParseVATCode(vatCode)
	return vat.Amount(price)
}
