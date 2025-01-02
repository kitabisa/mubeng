package awsurl

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *URL
		wantErr bool
	}{
		{
			name: "valid url without quotes",
			url:  "aws://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@ap-southeast-1",
			want: &URL{
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Region:          "ap-southeast-1",
			},
			wantErr: false,
		},
		{
			name: "valid url with quotes",
			url:  `aws://"AKIAIOSFODNN7EXAMPLE":"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"@"ap-southeast-1"`,
			want: &URL{
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Region:          "ap-southeast-1",
			},
			wantErr: false,
		},
		{
			name: "valid url with quotes and trailing slash",
			url:  `aws://"AKIAIOSFODNN7EXAMPLE":"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"@"ap-southeast-1"/`,
			want: &URL{
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Region:          "ap-southeast-1",
			},
			wantErr: false,
		},
		{
			name: "valid url with quotes and path",
			url:  `aws://"AKIAIOSFODNN7EXAMPLE":"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"@"ap-southeast-1"/ignored/path`,
			want: &URL{
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Region:          "ap-southeast-1",
			},
			wantErr: false,
		},
		{
			name: "mixed quotes",
			url:  `aws://"AKIAIOSFODNN7EXAMPLE":wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@ap-southeast-1`,
			want: &URL{
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Region:          "ap-southeast-1",
			},
			wantErr: false,
		},
		{
			name:    "invalid prefix",
			url:     "invalid://key:secret@region/",
			wantErr: true,
		},
		{
			name:    "missing @",
			url:     `aws://"key":"secret"`,
			wantErr: true,
		},
		{
			name:    "missing credentials separator",
			url:     "aws://keysecret@region/",
			wantErr: true,
		},
		{
			name:    "empty access key",
			url:     `aws://:"secret"@"region"`,
			wantErr: true,
		},
		{
			name:    "empty secret key",
			url:     `aws://"key":@"region"`,
			wantErr: true,
		},
		{
			name:    "empty region",
			url:     `aws://"key":"secret"@""`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.AccessKeyID != tt.want.AccessKeyID {
					t.Errorf("Parse() AccessKeyID = %v, want %v", got.AccessKeyID, tt.want.AccessKeyID)
				}
				if got.SecretAccessKey != tt.want.SecretAccessKey {
					t.Errorf("Parse() SecretAccessKey = %v, want %v", got.SecretAccessKey, tt.want.SecretAccessKey)
				}
				if got.Region != tt.want.Region {
					t.Errorf("Parse() Region = %v, want %v", got.Region, tt.want.Region)
				}
			}
		})
	}
}
