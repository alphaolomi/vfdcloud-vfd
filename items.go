package vfd

import (
	"github.com/vfdcloud/vfd/models"
)

type (
	ItemProcessResponse struct {
		ItemsXML          []*models.ITEM
		DISCOUNT          float64
		TOTALTAXEXCLUSIVE float64
		TOTALTAXINCLUSIVE float64
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
	var itemsXML []*models.ITEM
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
		itemsXML = append(itemsXML, itemXML)
		netAmount := NetAmount(item.TaxCode, itemAmount)
		TOTALTAXEXCLUSIVE += netAmount
		TOTALTAXINCLUSIVE += itemAmount
	}
	return &ItemProcessResponse{
		ItemsXML:          itemsXML,
		DISCOUNT:          DISCOUNT,
		TOTALTAXEXCLUSIVE: TOTALTAXEXCLUSIVE,
		TOTALTAXINCLUSIVE: TOTALTAXINCLUSIVE,
	}
}
