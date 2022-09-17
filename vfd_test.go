package vfd

import (
	"github.com/vfdcloud/vfdcloud/pkg/env"
	"testing"
)

func TestGetRequestURL(t *testing.T) {
	t.Parallel()
	type args struct {
		e   env.Env
		req RequestType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "staging registration",
			args: args{
				e:   env.STAGING,
				req: RequestTypeRegister,
			},
			want: "https://virtual.tra.go.tz/efdmsRctApi/api/vfdRegReq",
		},
		{
			name: "staging token",
			args: args{
				e:   env.STAGING,
				req: RequestTypeToken,
			},
			want: "https://virtual.tra.go.tz/efdmsRctApi/vfdtoken",
		},
		{
			name: "staging receipt",
			args: args{
				e:   env.STAGING,
				req: RequestTypeReceipt,
			},
			want: "https://virtual.tra.go.tz/efdmsRctApi/api/efdmsRctInfo",
		},
		{
			name: "staging report",
			args: args{
				e:   env.STAGING,
				req: RequestTypeReport,
			},
			want: "https://virtual.tra.go.tz/efdmsRctApi/api/efdmszreport",
		},
		{
			name: "production registration",
			args: args{
				e:   env.PRODUCTION,
				req: RequestTypeRegister,
			},
			want: "https://vfd.tra.go.tz/api/vfdRegReq",
		},
		{
			name: "production token",
			args: args{
				e:   env.PRODUCTION,
				req: RequestTypeToken,
			},
			want: "https://vfd.tra.go.tz/vfdtoken",
		},
		{
			name: "production receipt",
			args: args{
				e:   env.PRODUCTION,
				req: RequestTypeReceipt,
			},
			want: "https://vfd.tra.go.tz/api/efdmsRctInfo",
		},
		{
			name: "production report",
			args: args{
				e:   env.PRODUCTION,
				req: RequestTypeReport,
			},
			want: "https://vfd.tra.go.tz/api/efdmszreport",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GetRequestURL(tt.args.e, tt.args.req); got != tt.want {
				t.Errorf("GetRequestURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
