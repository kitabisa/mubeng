package awsurl

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// URL represents a parsed AWS URL containing credentials and region
type URL struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

// Parse parses an AWS URL in the format:
// - aws://ACCESS_KEY_ID:SECRET_ACCESS_KEY@region
//
// Supports optional quotes around credentials and region.
func Parse(rawURL string) (*URL, error) {
	// Check prefix
	if !IsURL(rawURL) {
		return nil, fmt.Errorf("%w: must start with aws://", ErrInvalidAWSURL)
	}

	// Remove prefix
	urlWithoutPrefix := strings.TrimPrefix(rawURL, "aws://")

	// Split into credentials and region parts
	parts := strings.Split(urlWithoutPrefix, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: missing @ separator", ErrInvalidAWSURL)
	}

	credentials := parts[0]

	// Extract region and remove any path after it
	region := parts[1]
	if idx := strings.Index(region, "/"); idx != -1 {
		region = region[:idx]
	}

	// Parse credentials
	var accessKeyID, secretAccessKey string
	credParts := strings.Split(credentials, ":")
	if len(credParts) != 2 {
		return nil, fmt.Errorf("%w: %w", ErrInvalidAWSURL, ErrInvalidCredentialsFormat)
	}

	// Remove quotes and trim spaces from credentials
	accessKeyID = strings.TrimSpace(unquote(credParts[0]))
	secretAccessKey = strings.TrimSpace(unquote(credParts[1]))

	// Remove quotes and trim spaces from region
	region = strings.TrimSpace(unquote(region))

	// Validate parts
	if accessKeyID == "" {
		return nil, fmt.Errorf("%w: ACCESS_KEY_ID cannot be empty", ErrInvalidAWSURL)
	}
	if secretAccessKey == "" {
		return nil, fmt.Errorf("%w: SECRET_ACCESS_KEY cannot be empty", ErrInvalidAWSURL)
	}
	if region == "" {
		return nil, fmt.Errorf("%w: region cannot be empty", ErrInvalidAWSURL)
	}

	return &URL{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
	}, nil
}

// IsURL checks if a string is an AWS URL
func IsURL(s string) bool {
	return strings.HasPrefix(s, "aws://")
}

// IsValidURL checks if a string is a valid AWS URL
func IsValidURL(s string) bool {
	_, err := Parse(s)
	return err == nil
}

// String returns the string representation of the URL
func (u *URL) String() string {
	return fmt.Sprintf("aws://%s:%s@%s", u.AccessKeyID, u.SecretAccessKey, u.Region)
}

// Credentials returns AWS credentials for the URL
func (u *URL) Credentials(session string) (aws.Credentials, error) {
	creds := credentials.NewStaticCredentialsProvider(u.AccessKeyID, u.SecretAccessKey, session)

	return creds.Retrieve(context.TODO())
}
