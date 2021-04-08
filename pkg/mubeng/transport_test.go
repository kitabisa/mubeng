package mubeng

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"

	"github.com/go-test/deep"
	"golang.org/x/net/proxy"
)

func TestTransport(t *testing.T) {
	type args struct {
		p string
	}

	failProxy := "gopher://localhost:70"
	httpProxy := "http://localhost:80"
	socks5Proxy := "socks5://localhost:3128"

	httpURL, err := url.Parse(httpProxy)
	if err != nil {
		t.Fatal(err)
	}

	socks5URL, err := url.Parse(socks5Proxy)
	if err != nil {
		t.Fatal(err)
	}

	dialer, err := proxy.SOCKS5("tcp", socks5URL.Host, nil, proxy.Direct)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		wantTr  *http.Transport
		wantErr bool
	}{
		{
			name: "Switch-transport to HTTP",
			args: args{
				p: httpProxy,
			},
			wantTr: &http.Transport{
				Proxy:             http.ProxyURL(httpURL),
				DisableKeepAlives: true,
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			},
			wantErr: false,
		},
		{
			name: "Switch-transport to SOCKSv5",
			args: args{
				p: socks5Proxy,
			},
			wantTr: &http.Transport{
				Dial:              dialer.Dial,
				DisableKeepAlives: true,
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			},
			wantErr: false,
		},
		{
			name: "Unsupported proxy proxy protocol scheme",
			args: args{
				p: failProxy,
			},
			wantTr:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTr, err := Transport(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(gotTr, tt.wantTr); diff != nil {
				t.Error(diff)
			}
		})
	}
}
