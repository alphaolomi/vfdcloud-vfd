package vfd

import "testing"

func TestParsePayment(t *testing.T) {
	t.Parallel()
	type test struct {
		name  string
		value any
		want  PaymentType
	}

	tests := []test{
		{
			name:  "TestParsePayment",
			value: 1,
			want:  CashPaymentType,
		},
		{
			name:  "TestParsePayment",
			value: "cash",
			want:  CashPaymentType,
		},
		{
			name:  "TestParsePayment",
			value: "CHEQUE",
			want:  ChequePaymentType,
		},
		{
			name:  "TestParsePayment",
			value: 3,
			want:  CreditCardPaymentType,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ParsePayment(tt.value); got != tt.want {
				t.Errorf("ParsePayment() = %v, want %v", got, tt.want)
			}
		})
	}
}
