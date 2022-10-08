package vfd

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestReportBytes(t *testing.T) {
	type args struct {
		privateKey *rsa.PrivateKey
		params     *ReportParams
		address    Address
		vats       []VATTOTAL
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
			name: "",
			args: args{
				privateKey: &rsa.PrivateKey{},
				params: &ReportParams{
					Date:             "2022-01-01",
					Time:             "12:00:00",
					VRN:              "123456789",
					TIN:              "T02FDC124",
					TaxOffice:        "London, UK",
					RegistrationID:   "123456789",
					ZNumber:          "20220101",
					EFDSerial:        "3CF241",
					RegistrationDate: "2021-01-01",
				},
				address: Address{
					Name:    "John Doe",
					Street:  "123 Main Street",
					Mobile:  "123456789",
					City:    "London",
					Country: "UK",
				},
				vats: []VATTOTAL{
					{
						ID:     StandardVATID,
						Rate:   StandardVATRATE,
						Amount: 20000000000,
					},
				},
				payments: []Payment{
					{
						Type:   CashPaymentType,
						Amount: 200000000000,
					},
				},
				totals: ReportTotals{
					DailyTotalAmount: 20000000000000,
					Gross:            20000000000000000,
					Corrections:      0,
					Discounts:        0,
					Surcharges:       0,
					TicketsVoid:      0,
					TicketsVoidTotal: 0,
					TicketsFiscal:    20000,
					TicketsNonFiscal: 0,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bitSize := 4096
			key, err := rsa.GenerateKey(rand.Reader, bitSize)
			if err != nil {
				t.Errorf("failed to generate private key: %s", err)
			}
			got, err := ReportBytes(key, tt.args.params, tt.args.address, tt.args.vats, tt.args.payments, tt.args.totals)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReportBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// PRINT THE BYTES AS STRING
			t.Logf("ReportBytes() = %v", string(got))
		})
	}
}
