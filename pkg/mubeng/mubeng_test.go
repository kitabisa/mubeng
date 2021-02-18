package mubeng

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestProxy_New(t *testing.T) {
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
		want1  *http.Request
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
			want1: req,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxy := &Proxy{
				Address:   tt.fields.Address,
				Transport: tt.fields.Transport,
			}
			got, got1 := proxy.New(tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Proxy.New() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Proxy.New() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
