package vfd

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestReportPayloadBytes(t *testing.T) {

	// generate a private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		privateKey *rsa.PrivateKey
		params     *ReportParams
		address    Address
		vats       []VatTotal
		payments   []Payment
		totals     ReportTotals
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				privateKey: privateKey,
				params: &ReportParams{
					Date:             "2020-09-04",
					Time:             "17:08:04",
					VRN:              "123456789",
					TIN:              "123456789",
					TaxOffice:        "123456789",
					RegistrationID:   "123456789",
					ZNumber:          "123456789",
					EFDSerial:        "123456789",
					RegistrationDate: "2020-09-04",
					User:             "123456789",
				},
				address: Address{
					Name:    "Example Name User",
					Street:  "26 Road, OHIO 12345 STREET",
					Mobile:  "+123456789",
					City:    "Miami",
					Country: "Las Vegas",
				},
				vats: []VatTotal{
					{
						ID:     "A",
						Rate:   0.18,
						Amount: 15000,
					},
					{
						ID:     "B",
						Rate:   0.0,
						Amount: 15000,
					},
					{
						ID:     "C",
						Rate:   0,
						Amount: 15000,
					},
				},
				payments: []Payment{
					{
						Type:   CashPaymentType,
						Amount: 145.90,
					},
					{
						Type:   CreditCardPaymentType,
						Amount: 6146.90,
					},
					{
						Type:   ElectronicPaymentType,
						Amount: 50000.90,
					},
				},
				totals: ReportTotals{
					DailyTotalAmount: 101010,
					Gross:            627272.9,
					Corrections:      272727.27,
					Discounts:        670.90,
					Surcharges:       6454.4,
					TicketsVoid:      13,
					TicketsVoidTotal: 151671.5,
					TicketsFiscal:    700,
					TicketsNonFiscal: 89,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReportPayloadBytes(tt.args.privateKey, tt.args.params, tt.args.address, tt.args.vats, tt.args.payments, tt.args.totals)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReportPayloadBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ReportPayloadBytes() = %v", string(got))
		})
	}
}
