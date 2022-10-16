package vfd_test

import (
	"github.com/vfdcloud/vfd"
	"testing"
)

func TestNetAmountAndTaxAmount(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name    string
		vatCode int64
		amount  float64
	}

	testCases := []testCase{
		{"taxable items", vfd.TaxableItemCode, 10000},
		{"non-taxable items", vfd.NonTaxableItemCode, 10000},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			netAmount := vfd.NetAmount(tc.vatCode, tc.amount)
			taxAmount := vfd.ValueAddedTaxAmount(tc.vatCode, tc.amount)
			totalAmount := netAmount + taxAmount
			if totalAmount != tc.amount {
				t.Errorf("totalAmount %.2f != amount %.2f", totalAmount, tc.amount)
			}
			t.Logf("netAmount: %.2f, taxAmount: %.2f, totalAmount: %.2f",
				netAmount, taxAmount, totalAmount)
		})
	}
}

func TestReportTaxRateID(t *testing.T) {
	type args struct {
		taxCode int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "taxable items",
			args: args{taxCode: vfd.TaxableItemCode},
			want: "A-18.00",
		},
		{
			name: "non-taxable items",
			args: args{taxCode: vfd.NonTaxableItemCode},
			want: "C-0.00",
		},
		{
			name: "standard vat",
			args: args{taxCode: vfd.StandardVATCODE},
			want: "A-18.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vfd.ReportTaxRateID(tt.args.taxCode); got != tt.want {
				t.Errorf("ReportTaxRateID() = %v, want %v", got, tt.want)
			}
		})
	}
}
