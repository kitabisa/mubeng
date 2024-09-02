package mubeng

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/go-test/deep"
)

func TestProxyNew(t *testing.T) {
	type fields struct {
		Address   string
		Transport *http.Transport
	}

	type args struct {
		req *http.Request
	}

	address := "http://localhost:3128"

	proxyURL, err := url.Parse(address)
	if err != nil {
		t.Fatal(err)
	}

	tr := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *http.Client
		want1  error
	}{
		{
			name: "New Proxy",
			fields: fields{
				Address:   address,
				Transport: tr,
			},
			args: args{
				req: req,
			},
			want: &http.Client{
				Transport: tr,
			},
			want1: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxy := &Proxy{
				Address:   tt.fields.Address,
				Transport: tt.fields.Transport,
			}
			got, got1 := proxy.New(tt.args.req)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Error(diff)
			}
			if diff1 := deep.Equal(got1, tt.want1); diff1 != nil {
				t.Error(diff1)
			}
		})
	}
}
