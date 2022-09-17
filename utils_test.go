package vfd

import (
	"github.com/vfdcloud/vfdcloud/pkg/env"
	"testing"
)

func TestReceiptLink(t *testing.T) {
	type args struct {
		e                         env.Env
		receiptVerificationNumber string
		receiptVerificationTime   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "dev receipt link",
			args: args{
				e:                         env.STAGING,
				receiptVerificationNumber: "123456789",
				receiptVerificationTime:   "12:00:00",
			},
			want: "https://virtual.tra.go.tz/efdmsRctVerify/123456789_120000",
		},
		{
			name: "prod receipt link",
			args: args{
				e:                         env.PRODUCTION,
				receiptVerificationNumber: "123456789",
				receiptVerificationTime:   "12:00:00",
			},
			want: "https://verify.tra.go.tz/123456789_120000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReceiptLink(tt.args.e, tt.args.receiptVerificationNumber, tt.args.receiptVerificationTime); got != tt.want {
				t.Errorf("ReceiptLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
