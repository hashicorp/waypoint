package cli

import "testing"

func Test_addUrlScheme(t *testing.T) {
	hostname := "hostname.com"
	httpsSchemeStr := string(httpsScheme + "://")
	httpSchemeStr := string(httpScheme + "://")

	type args struct {
		urlString string
		scheme    UrlScheme
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid no scheme URL to add https://",
			args: args{
				urlString: hostname,
				scheme:    httpsScheme,
			},
			want:    string(httpsSchemeStr + hostname),
			wantErr: false,
		},
		{
			name: "Valid https URL to add https://",
			args: args{
				urlString: string(httpsSchemeStr + hostname),
				scheme:    httpsScheme,
			},
			want:    string(httpsSchemeStr + hostname),
			wantErr: false,
		},
		{
			name: "Valid http URL to add https://",
			args: args{
				urlString: string(httpSchemeStr + hostname),
				scheme:    httpsScheme,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Empty URL to add https://",
			args: args{
				urlString: "",
				scheme:    httpsScheme,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addUrlScheme(tt.args.urlString, tt.args.scheme)
			if (err != nil) != tt.wantErr {
				t.Errorf("addUrlScheme() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("addUrlScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}
