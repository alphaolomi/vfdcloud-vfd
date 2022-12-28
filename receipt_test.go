package vfd

import (
	"encoding/xml"
	"fmt"
	"github.com/vfdcloud/vfd/internal/models"
	"reflect"
	"testing"
)

func TestProcessItems(t *testing.T) {
	type args struct {
		items []Item
	}
	type result struct {
		ItemsAmount []float64
		VAT         []*models.VATTOTAL
		Totals      models.TOTALS
	}
	tests := []struct {
		name string
		args args
		want *result
	}{
		{
			name: "TestProcessItems",
			args: args{
				items: []Item{
					{
						ID:          "1",
						Description: "Item 1",
						TaxCode:     TaxableItemCode,
						Quantity:    5,
						Price:       2000,
						Discount:    1000,
					},
				},
			},
			want: &result{
				ItemsAmount: []float64{10000},
				VAT: []*models.VATTOTAL{
					{
						XMLName:    xml.Name{},
						Text:       "",
						VATRATE:    "A",
						NETTAMOUNT: "4237.29",
						TAXAMOUNT:  "762.71",
					},
				},
				Totals: models.TOTALS{
					TOTALTAXEXCL: 4237.29,
					TOTALTAXINCL: 5000.00,
					DISCOUNT:     5000.00,
				},
			},
		},
		{
			name: "TestProcessItems",
			args: args{
				items: []Item{
					{
						ID:          "1",
						Description: "Item 1",
						TaxCode:     TaxableItemCode,
						Quantity:    5,
						Price:       1000,
						Discount:    0,
					},
				},
			},
			want: &result{
				ItemsAmount: []float64{5000},
				VAT: []*models.VATTOTAL{
					{
						XMLName:    xml.Name{},
						Text:       "",
						VATRATE:    "A",
						NETTAMOUNT: "4237.29",
						TAXAMOUNT:  "762.71",
					},
				},
				Totals: models.TOTALS{
					TOTALTAXEXCL: 4237.29,
					TOTALTAXINCL: 5000.00,
					DISCOUNT:     0.00,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessItems(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				items := got.ITEMS
				for i, item := range items {
					wantItemAmount := tt.want.ItemsAmount[i]
					message := fmt.Sprintf("[GOT]: Item %d: id: %s, quantity: %.2f, amount: %.2f [EXPECTED]: %.2f\n", i, item.ID, item.QTY, item.AMT, wantItemAmount)
					if item.AMT != wantItemAmount {
						t.Errorf("ProcessItems(): %s", message)
					}
					t.Logf("ProcessItems(): %s", message)
				}

				// Comparing TOTALS
				func(got, want models.TOTALS) {
					if got.TOTALTAXEXCL != want.TOTALTAXEXCL {
						t.Errorf("TOTALTAXEXCL: got %.2f, want %.2f", got.TOTALTAXEXCL, want.TOTALTAXEXCL)
					}
					if got.TOTALTAXINCL != want.TOTALTAXINCL {
						t.Errorf("TOTALTAXINCL: got %.2f, want %.2f", got.TOTALTAXINCL, want.TOTALTAXINCL)
					}
					if got.DISCOUNT != want.DISCOUNT {
						t.Errorf("DISCOUNT: got %.2f, want %.2f", got.DISCOUNT, want.DISCOUNT)
					}
					t.Logf("TOTALS: [GOT]: TAXEXCL: %.2f, TAXINCL: %.2f, DISCOUNT: %.2f [EXPECTED]: TAXEXCL: %.2f, TAXINCL: %.2f, DISCOUNT: %.2f ",
						got.TOTALTAXEXCL, got.TOTALTAXINCL, got.DISCOUNT, want.TOTALTAXEXCL, want.TOTALTAXINCL, want.DISCOUNT)
				}(got.TOTALS, tt.want.Totals)

				// Comparing VATS
				func(got, want []*models.VATTOTAL) {
					for i, v := range got {
						if v.TAXAMOUNT != want[i].TAXAMOUNT {
							t.Errorf("VATRATE[%s]: TAXAMOUNT: got %s, want %s", v.VATRATE, v.TAXAMOUNT, want[i].TAXAMOUNT)
						}
						if v.NETTAMOUNT != want[i].NETTAMOUNT {
							t.Errorf("VATRATE[%s].NETTAMOUNT: got %s, want %s", v.VATRATE, v.NETTAMOUNT, want[i].NETTAMOUNT)
						}
						t.Logf("VATRATE[%s]: [GOT]: TAXAMOUNT: %s, NETTAMOUNT: %s [EXPECTED]: TAXAMOUNT: %s, NETTAMOUNT: %s ",
							v.VATRATE, v.TAXAMOUNT, v.NETTAMOUNT, want[i].TAXAMOUNT, want[i].NETTAMOUNT)
					}

				}(got.VATTOTALS, tt.want.VAT)

			}
		})
	}
}
