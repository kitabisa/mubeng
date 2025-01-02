package awsurl

import "errors"

var (
	ErrInvalidAWSURL            = errors.New("invalid AWS URL")
	ErrInvalidCredentialsFormat = errors.New("credentials must be in format ACCESS_KEY_ID:SECRET_ACCESS_KEY")
)
