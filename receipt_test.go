package vfd

import (
	"encoding/xml"
	"github.com/vfdcloud/vfd/models"
	"strings"
	"testing"
)

func TestGenerateReceipt(t *testing.T) {
	type args struct {
		params   ReceiptParams
		customer Customer
		items    []Item
		payments []Payment
	}
	tests := []struct {
		name string
		args args
		want *models.RCT
	}{
		{
			name: "test",
			args: args{
				params: ReceiptParams{
					Date:           "1998-01-01",
					Time:           "14:00:00",
					TIN:            "123456789",
					RegistrationID: "123456789",
					EFDSerial:      "123456789",
					ReceiptNum:     "123456789",
					DailyCounter:   5,
					GlobalCounter:  100,
					ZNum:           "FARDSCCS",
					ReceiptVNum:    "GATFFARDA",
				},
				customer: Customer{
					Type:   NonCustomerID,
					ID:     "",
					Name:   "Pius Alfred",
					Mobile: "087716652442",
				},
				items: []Item{
					{
						ID:          "12345",
						Description: "Mlimani City Parking Charge",
						TaxCode:     1,
						Quantity:    1,
						Price:       10000,
						Discount:    0,
					},
				},
				payments: []Payment{
					{
						Type:   CashPaymentType,
						Amount: 10000,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateReceipt(tt.args.params, tt.args.customer, tt.args.items, tt.args.payments); got != nil {
				// Marshal the receipt (XML)
				b, err := xml.Marshal(got)
				if err != nil {
					t.Errorf("Marshal() error = %v", err)
					return
				}

				// creates a replacer that replaces all the characters that are not allowed in XML
				replacer := strings.NewReplacer(
					"<PAYMENT>", "",
					"</PAYMENT>", "",
					"<VATTOTAL>", "",
					"</VATTOTAL>", "")

				// replace all the characters that are not allowed in XML
				xmlString := replacer.Replace(string(b))

				t.Logf("GenerateReceipt() = %s", string(xmlString))
			}
		})
	}
}
