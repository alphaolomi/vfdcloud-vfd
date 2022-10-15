package vfd

import (
	"testing"
)

func TestProcessItems(t *testing.T) {
	type args struct {
		items []Item
	}
	tests := []struct {
		name string
		args args
		want *ItemProcessResponse
	}{
		{
			name: "ProcessItems",
			args: args{
				items: []Item{
					{
						ID:          "ID0001",
						Description: "Item 1",
						TaxCode:     StandardVATCODE,
						Quantity:    5,
						Price:       1000,
						Discount:    0,
					},
					{
						ID:          "ID0002",
						Description: "Item 2",
						TaxCode:     StandardVATCODE,
						Quantity:    4,
						Price:       2000,
					},
					{
						ID:          "ID0003",
						Description: "Item 3",
						TaxCode:     NonTaxableItemCode,
						Quantity:    3,
						Price:       3000,
						Discount:    0,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := ProcessItems(tt.args.items)

			msg := got.PrettyMessage()

			t.Log(msg)

		})
	}
}
