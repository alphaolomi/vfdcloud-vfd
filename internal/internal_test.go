package internal

import "testing"

func TestAppendEndpoint(t *testing.T) {
	type args struct {
		url       string
		endpoints []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				url: "https://www.google.com/",
				endpoints: []string{
					"/api/", "/v1/", "/users", "photos",
				},
			},
			want: "https://www.google.com/api/v1/users/photos",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendEndpoint(tt.args.url, tt.args.endpoints...); got != tt.want {
				t.Errorf("AppendEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
