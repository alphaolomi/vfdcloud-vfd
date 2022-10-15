package vfd

import (
	"github.com/vfdcloud/vfd/models"
)

type (
	ItemProcessResponse struct {
		ITEMS     []*models.ITEM
		VATTOTALS []*models.VATTOTAL
		TOTALS    models.TOTALS
	}
)

// ProcessItems processes the []Items in the submitted receipt request
// and create []*models.ITEM which is used to create the xml request also
// calculates the total discount, total tax exclusive and total tax inclusive
func ProcessItems(items []Item) *ItemProcessResponse {

	var (
		DISCOUNT          = 0.0
		TOTALTAXEXCLUSIVE = 0.0
		TOTALTAXINCLUSIVE = 0.0
	)
	// initialize map that will store the tax code and total amount of tax collected
	// over all items with the same tax code. The map keys are the tax codes.
	vatTotals := make(map[string]*models.VATTOTAL)
	var ITEMS []*models.ITEM
	for _, item := range items {
		itemAmount := item.Quantity * item.Price
		itemXML := &models.ITEM{
			ID:      item.ID,
			DESC:    item.Description,
			QTY:     item.Quantity,
			TAXCODE: item.TaxCode,
			AMT:     itemAmount,
		}
		DISCOUNT += item.Discount
		ITEMS = append(ITEMS, itemXML)
		NETAMOUNT := NetAmount(item.TaxCode, itemAmount)
		TOTALTAXEXCLUSIVE += NETAMOUNT
		TOTALTAXINCLUSIVE += itemAmount
		TAXAMOUNT := ValueAddedTaxAmount(item.TaxCode, itemAmount)

		vat := ParseTaxCode(item.TaxCode)
		vatID := vat.ID

		// check if the tax code is already in the map if not add it
		if _, ok := vatTotals[vatID]; !ok {
			vatTotals[vatID] = &models.VATTOTAL{
				VATRATE:    vatID,
				NETTAMOUNT: NETAMOUNT,
				TAXAMOUNT:  TAXAMOUNT,
			}
		} else {
			vatTotals[vatID].NETTAMOUNT += NETAMOUNT
			vatTotals[vatID].TAXAMOUNT += TAXAMOUNT
		}
	}

	VATTOTALS := make([]*models.VATTOTAL, 0)
	for _, v := range vatTotals {
		VATTOTALS = append(VATTOTALS, v)
	}
	TOTALS := models.TOTALS{
		TOTALTAXEXCL: TOTALTAXEXCLUSIVE,
		TOTALTAXINCL: TOTALTAXINCLUSIVE,
		DISCOUNT:     DISCOUNT,
	}
	return &ItemProcessResponse{
		ITEMS:     ITEMS,
		VATTOTALS: VATTOTALS,
		TOTALS:    TOTALS,
	}
}
