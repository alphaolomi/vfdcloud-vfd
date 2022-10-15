package vfd

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestReceiptBytes(t *testing.T) {
	// create a rsa private key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %s", err)
	}
	type args struct {
		privateKey *rsa.PrivateKey
		params     ReceiptParams
		customer   Customer
		items      []Item
		payments   []Payment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ReceiptBytes",
			args: args{
				privateKey: key,
				params: ReceiptParams{
					Date:           "2020-01-01",
					Time:           "12:00:00",
					TIN:            "123456789",
					RegistrationID: "123456789",
					EFDSerial:      "123456789",
					ReceiptNum:     "123456789",
					DailyCounter:   156,
					GlobalCounter:  2678,
					ZNum:           "202001012678",
					ReceiptVNum:    "120000156",
				},
				customer: Customer{
					Type:   NIDACustomerID,
					ID:     "526262",
					Name:   "Pius Alfred",
					Mobile: "255765992153",
				},
				items: []Item{
					{
						ID:          "ITEMOO1",
						Description: "Item 1",
						TaxCode:     StandardVATCODE,
						Quantity:    2,
						Price:       1500,
						Discount:    0,
					},
					{
						ID:          "ITEMOO2",
						Description: "Item 2",
						TaxCode:     StandardVATCODE,
						Quantity:    1,
						Price:       1000,
						Discount:    0,
					},
					{
						ID:          "ITEMOO3",
						Description: "Item 3",
						TaxCode:     3,
						Quantity:    1,
						Price:       500,
					},
				},
				payments: []Payment{
					{
						Type:   CashPaymentType,
						Amount: 1500,
					},
					{
						Type:   ElectronicPaymentType,
						Amount: 1000,
					},
					{
						Type:   CreditCardPaymentType,
						Amount: 3000,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReceiptBytes(tt.args.privateKey, tt.args.params, tt.args.customer, tt.args.items, tt.args.payments)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReceiptBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// print the receipt
			t.Logf("\n%s\n", got)
		})
	}
}
