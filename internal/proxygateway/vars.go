package proxygateway

// ValidRegionCodes is a map of valid AWS API Gateway region codes.
//
// This map is used to validate the region code string provided by the user.
// Data is taken from the official AWS documentation.
//
// ```bash
// curl -O https://ip-ranges.amazonaws.com/ip-ranges.json
// jq -r '.prefixes[] | select(.service=="API_GATEWAY") | "\"\(.region)\": true,"' < ip-ranges.json | sort -u
// ```
var ValidRegionCodes = map[string]bool{
	"af-south-1":     true,
	"ap-east-1":      true,
	"ap-northeast-1": true,
	"ap-northeast-2": true,
	"ap-northeast-3": true,
	"ap-south-1":     true,
	"ap-south-2":     true,
	"ap-southeast-1": true,
	"ap-southeast-2": true,
	"ap-southeast-3": true,
	"ap-southeast-4": true,
	"ap-southeast-5": true,
	"ap-southeast-7": true,
	"ca-central-1":   true,
	"ca-west-1":      true,
	"cn-north-1":     true,
	"cn-northwest-1": true,
	"eu-central-1":   true,
	"eu-central-2":   true,
	"eu-north-1":     true,
	"eu-south-1":     true,
	"eu-south-2":     true,
	"eu-west-1":      true,
	"eu-west-2":      true,
	"eu-west-3":      true,
	"il-central-1":   true,
	"me-central-1":   true,
	"me-south-1":     true,
	"mx-central-1":   true,
	"sa-east-1":      true,
	"us-east-1":      true,
	"us-east-2":      true,
	"us-gov-east-1":  true,
	"us-gov-west-1":  true,
	"us-west-1":      true,
	"us-west-2":      true,
}
