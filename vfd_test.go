package vfd

import (
	"github.com/vfdcloud/base"
	"testing"
)

func TestGetRequestURL(t *testing.T) {
	t.Parallel()
	type args struct {
		e   base.Env
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
				e:   base.StagingEnv,
				req: RequestTypeRegister,
			},
			want: "https://virtual.tra.go.tz/efdmsRctApi/api/vfdRegReq",
		},
		{
			name: "staging token",
			args: args{
				e:   base.StagingEnv,
				req: RequestTypeToken,
			},
			want: "https://virtual.tra.go.tz/efdmsRctApi/vfdtoken",
		},
		{
			name: "staging receipt",
			args: args{
				e:   base.StagingEnv,
				req: RequestTypeReceipt,
			},
			want: "https://virtual.tra.go.tz/efdmsRctApi/api/efdmsRctInfo",
		},
		{
			name: "staging report",
			args: args{
				e:   base.StagingEnv,
				req: RequestTypeReport,
			},
			want: "https://virtual.tra.go.tz/efdmsRctApi/api/efdmszreport",
		},
		{
			name: "production registration",
			args: args{
				e:   base.ProdEnv,
				req: RequestTypeRegister,
			},
			want: "https://vfd.tra.go.tz/api/vfdRegReq",
		},
		{
			name: "production token",
			args: args{
				e:   base.ProdEnv,
				req: RequestTypeToken,
			},
			want: "https://vfd.tra.go.tz/vfdtoken",
		},
		{
			name: "production receipt",
			args: args{
				e:   base.ProdEnv,
				req: RequestTypeReceipt,
			},
			want: "https://vfd.tra.go.tz/api/efdmsRctInfo",
		},
		{
			name: "production report",
			args: args{
				e:   base.ProdEnv,
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
