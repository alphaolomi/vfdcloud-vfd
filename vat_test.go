package vfd

import (
	"math"
	"testing"
)

func TestNetPrice(t *testing.T) {
	type args struct {
		taxCode int64
		price   float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "StandardVAT",
			args: args{
				taxCode: StandardVATCODE,
				price:   1000,
			},
		},
		{
			name: "StandardVAT",
			args: args{
				taxCode: StandardVATCODE,
				price:   5000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NetAmount(tt.args.taxCode, tt.args.price)
			vat := parseTaxCode(tt.args.taxCode)
			rate := vat.Percentage / 100
			total := got + (got * rate)
			total = math.Floor(total*100) / 100
			if total != tt.args.price {
				t.Errorf("NetAmount() = %v, want %v", total, tt.args.price)
			}
		})
	}
}

func TestValueAddedTaxAmount(t *testing.T) {
	type args struct {
		taxCode int64
		price   float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "StandardVAT",
			args: args{
				taxCode: StandardVATCODE,
				price:   1000,
			},
		},
		{
			name: "StandardVAT",
			args: args{
				taxCode: StandardVATCODE,
				price:   5000,
			},
		},
		{
			name: "StandardVAT",
			args: args{
				taxCode: StandardVATCODE,
				price:   200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tax := ValueAddedTaxAmount(tt.args.taxCode, tt.args.price)
			netAmount := NetAmount(tt.args.taxCode, tt.args.price)
			total := netAmount + tax
			total = math.Floor(total*100) / 100
			t.Logf("tax: %v, netAmount: %v, total: %v", tax, netAmount, total)
			if total != tt.args.price {
				t.Errorf("ValueAddedTaxAmount() = %v, want %v", total, tt.args.price)
			}
		})
	}
}
