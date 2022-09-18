package vfd

import (
	"github.com/vfdcloud/base"
	"testing"
)

func TestRequestURL(t *testing.T) {
	t.Parallel()
	type args struct {
		e      base.Env
		action Action
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "production registration url",
			args: args{
				e:      base.ProdEnv,
				action: RegisterClientAction,
			},
			want: "https://vfd.tra.go.tz/api/vfdRegReq",
		},
		{
			name: "staging verification url",
			args: args{
				e:      base.TestEnv,
				action: ReceiptVerificationAction,
			},
			want: "https://virtual.tra.go.tz/efdmsRctVerify/",
		},
	}
	for _, tt := range tests {
		tt := tt // Parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := RequestURL(tt.args.e, tt.args.action); got != tt.want {
				t.Errorf("RequestURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
