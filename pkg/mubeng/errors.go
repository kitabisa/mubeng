package mubeng

import "fmt"

var (
	ErrUnsupportedProxyProtocolScheme   = fmt.Errorf("unsupported proxy protocol scheme")
	ErrSwitchTransportAWSProtocolScheme = fmt.Errorf("could not switch transport for AWS protocol scheme")
)
