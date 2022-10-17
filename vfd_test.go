package vfd

import "testing"

func TestParsePayment(t *testing.T) {
	type args struct {
		value any
	}

	type test struct {
		name string
		args args
		want PaymentType
	}

	tests := []test{
		{
			name: "TestParsePayment",
			args: args{
				value: 1,
			},
			want: CashPaymentType,
		},
		{
			name: "TestParsePayment",
			args: args{
				value: "cash",
			},
			want: CashPaymentType,
		},
		{
			name: "TestParsePayment",
			args: args{
				value: "CHEQUE",
			},
			want: ChequePaymentType,
		},
		{
			name: "TestParsePayment",
			args: args{
				value: 3,
			},
			want: CreditCardPaymentType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePayment(tt.args.value); got != tt.want {
				t.Errorf("ParsePayment() = %v, want %v", got, tt.want)
			}
		})
	}
}
