package proxygateway

import (
	"context"
	"fmt"
	"io"
	"strings"

	"net/http"
	"net/url"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/kitabisa/mubeng/pkg/helper/awsurl"
)

// ProxyGateway contains information to interact with AWS API Gateway.
type ProxyGateway struct {
	baseURL       string
	parsedBaseURL *url.URL
	endpoint      *url.URL

	accessKeyID     string
	secretAccessKey string

	client     *apigateway.Client
	httpClient *http.Client

	region string
	apiID  *string
}

// New creates a new ProxyGateway instance
func New(ctx context.Context, accessKeyID, secretAccessKey, region string) (*ProxyGateway, error) {
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")
	if _, err := creds.Retrieve(ctx); err != nil {
		return nil, fmt.Errorf("could not retrieve credentials: %v", err)
	}

	if !ValidRegionCodes[region] {
		return nil, fmt.Errorf("invalid %q region", region)
	}

	pg := &ProxyGateway{
		baseURL:         "",
		region:          region,
		apiID:           nil,
		endpoint:        nil,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		httpClient:      new(http.Client),
	}

	if err := pg.createClient(ctx, creds); err != nil {
		return nil, err
	}

	return pg, nil
}

// NewFromURL creates a new ProxyGateway instance from an AWS URL
func NewFromURL(ctx context.Context, s string) (*ProxyGateway, error) {
	u, err := awsurl.Parse(s)
	if err != nil {
		return nil, err
	}

	return New(ctx, u.AccessKeyID, u.SecretAccessKey, u.Region)
}

// createClient creates a new API Gateway client
func (pg *ProxyGateway) createClient(ctx context.Context, credentials credentials.StaticCredentialsProvider) error {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(pg.region),
		config.WithCredentialsProvider(credentials),
	)

	if err != nil {
		return fmt.Errorf("unable to load SDK config for region %q: %v", pg.region, err)
	}

	pg.client = apigateway.NewFromConfig(cfg)

	return nil
}

// storeEndpoint stores the API endpoint URL
func (pg *ProxyGateway) storeEndpoint(region, apiID string) error {
	endpoint := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/%s", apiID, region, StageName)
	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("could not parse endpoint URL: %v", err)
	}

	pg.endpoint = parsedEndpoint

	return nil
}

// SetBaseURL sets the base URL for the gateway
func (pg *ProxyGateway) SetBaseURL(u string) error {
	var err error

	pg.baseURL, pg.parsedBaseURL, err = GetBaseURL(u)
	if err != nil {
		return err
	}

	return nil
}

// SetHTTPClient sets the HTTP client to be used for requests
func (pg *ProxyGateway) SetHTTPClient(client *http.Client) {
	pg.httpClient = client
}

// Start creates the API Gateway and deploys it
func (pg *ProxyGateway) Start(ctx context.Context) error {
	if pg.parsedBaseURL == nil {
		return fmt.Errorf("base URL not set")
	}

	// List existing APIs
	apis, err := pg.client.GetRestApis(ctx, new(apigateway.GetRestApisInput))
	if err != nil {
		return fmt.Errorf("could not list REST APIs: %v", err)
	}

	// Look for existing API with matching tags
	var existingAPI *types.RestApi
	for _, api := range apis.Items {
		if api.Tags["baseURL"] == pg.baseURL && api.Tags["region"] == pg.region {
			existingAPI = &api
			break
		}
	}

	if existingAPI != nil {
		pg.apiID = existingAPI.Id
	} else {
		// Create new API if none exists
		api, err := pg.client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:        aws.String(fmt.Sprintf("mubeng-proxy-gateway-%s", pg.region)),
			Description: aws.String(fmt.Sprintf("mubeng Proxy Gateway (%s)", pg.baseURL)),
			Tags: map[string]string{
				"baseURL": pg.baseURL,
				"region":  pg.region,
			},
			EndpointConfiguration: &types.EndpointConfiguration{
				Types: []types.EndpointType{types.EndpointTypeRegional},
			},
		})

		if err != nil {
			return fmt.Errorf("could not create REST API: %v", err)
		}

		pg.apiID = api.Id
	}

	// If API already exists, just store endpoint and return
	if existingAPI != nil {
		return pg.storeEndpoint(pg.region, *pg.apiID)
	}

	// Get root resource ID
	resources, err := pg.client.GetResources(ctx, &apigateway.GetResourcesInput{
		RestApiId: pg.apiID,
	})
	if err != nil {
		return fmt.Errorf("could not get resources: %v", err)
	}

	var rootResourceID string
	for _, resource := range resources.Items {
		if *resource.Path == "/" {
			rootResourceID = *resource.Id
			break
		}
	}

	// Create proxy resource
	proxyResource, err := pg.client.CreateResource(ctx, &apigateway.CreateResourceInput{
		RestApiId: pg.apiID,
		ParentId:  aws.String(rootResourceID),
		PathPart:  aws.String("{proxy+}"),
	})
	if err != nil {
		return fmt.Errorf("could not create proxy resource: %v", err)
	}

	_, err = pg.client.PutMethod(ctx, &apigateway.PutMethodInput{
		RestApiId:         pg.apiID,
		ResourceId:        aws.String(rootResourceID),
		HttpMethod:        aws.String("ANY"),
		AuthorizationType: aws.String("NONE"),
		RequestParameters: map[string]bool{
			"method.request.path.proxy":                    true,
			"method.request.header.X-Mubeng-Forwarded-For": false,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create ANY method: %v", err)
	}

	_, err = pg.client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
		RestApiId:             pg.apiID,
		ResourceId:            aws.String(rootResourceID),
		HttpMethod:            aws.String("ANY"),
		IntegrationHttpMethod: aws.String("ANY"),
		Uri:                   aws.String(fmt.Sprintf("%s/{proxy}", pg.baseURL)),
		Type:                  types.IntegrationTypeHttpProxy,
		RequestParameters: map[string]string{
			"integration.request.path.proxy":             "method.request.path.proxy",
			"integration.request.header.X-Forwarded-For": "method.request.header.X-Mubeng-Forwarded-For",
		},
	})
	if err != nil {
		return fmt.Errorf("could not create integration: %v", err)
	}

	_, err = pg.client.PutMethod(ctx, &apigateway.PutMethodInput{
		RestApiId:         pg.apiID,
		ResourceId:        proxyResource.Id,
		HttpMethod:        aws.String("ANY"),
		AuthorizationType: aws.String("NONE"),
		RequestParameters: map[string]bool{
			"method.request.path.proxy":                    true,
			"method.request.header.X-Mubeng-Forwarded-For": false,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create ANY method: %v", err)
	}

	_, err = pg.client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
		RestApiId:             pg.apiID,
		ResourceId:            proxyResource.Id,
		HttpMethod:            aws.String("ANY"),
		IntegrationHttpMethod: aws.String("ANY"),
		Uri:                   aws.String(fmt.Sprintf("%s/{proxy}", pg.baseURL)),
		Type:                  types.IntegrationTypeHttpProxy,
		RequestParameters: map[string]string{
			"integration.request.path.proxy":             "method.request.path.proxy",
			"integration.request.header.X-Forwarded-For": "method.request.header.X-Mubeng-Forwarded-For",
		},
	})
	if err != nil {
		return fmt.Errorf("could not create integration: %v", err)
	}

	// Deploy API
	_, err = pg.client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
		RestApiId: pg.apiID,
		StageName: aws.String(StageName),
	})
	if err != nil {
		return fmt.Errorf("could not create deployment: %v", err)
	}

	// Store endpoint
	return pg.storeEndpoint(pg.region, *pg.apiID)
}

// GetEndpoint returns the API Gateway endpoint
func (pg *ProxyGateway) GetEndpoint() *url.URL {
	return pg.endpoint
}

// Send sends an HTTP request through the gateway
func (pg *ProxyGateway) Send(method, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	if pg.endpoint == nil {
		return nil, fmt.Errorf("no available endpoint, make sure gateway is started")
	}

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		parsedURL, err := url.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("invalid path URL: %v", err)
		}

		path = parsedURL.Path
	}

	path = strings.TrimPrefix(path, "/")

	pg.endpoint.Path = filepath.Join(pg.endpoint.Path, path)
	url := pg.endpoint.String()

	// Create request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}

	// Copy headers
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Del("Host")

	if x := req.Header.Get("X-Forwarded-For"); x != "" {
		req.Header.Add("X-Mubeng-Forwarded-For", x)
		req.Header.Del("X-Forwarded-For")
	}

	// Send request
	resp, err := pg.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed through region %s: %v", pg.region, err)
	}

	return resp, nil
}

// SendRequest sends an [http.Request] through the gateway
func (pg *ProxyGateway) SendRequest(req *http.Request) (*http.Response, error) {
	headers := req.Header.Clone()

	return pg.Send(req.Method, req.URL.String(), req.Body, headers)
}

// Close cleans up the gateway resources
func (pg *ProxyGateway) Close(ctx context.Context) error {
	_, err := pg.client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
		RestApiId: pg.apiID,
	})

	return err
}

// Clean deletes all resources created by the gateway
func (pg *ProxyGateway) Clean(ctx context.Context) error {
	apis, err := pg.client.GetRestApis(ctx, new(apigateway.GetRestApisInput))
	if err != nil {
		return fmt.Errorf("could not list REST APIs: %v", err)
	}

	for _, api := range apis.Items {
		if strings.HasPrefix(*api.Name, "mubeng-proxy-gateway") {
			_, err := pg.client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
				RestApiId: api.Id,
			})
			if err != nil {
				return fmt.Errorf("could not delete API %s: %v", *api.Name, err)
			}
		}
	}

	return nil
}
