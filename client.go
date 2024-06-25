package zhipu

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt"
)

const (
	envAPIKey  = "ZHIPUAI_API_KEY"
	envBaseURL = "ZHIPUAI_BASE_URL"

	defaultBaseURL = "https://open.bigmodel.cn/api/paas/v4"
)

var (
	// ErrAPIKeyMissing is the error when the api key is missing
	ErrAPIKeyMissing = errors.New("zhipu: api key is missing")
	// ErrAPIKeyMalformed is the error when the api key is malformed
	ErrAPIKeyMalformed = errors.New("zhipu: api key is malformed")
)

type clientOptions struct {
	baseURL string
	apiKey  string
	debug   bool
}

// ClientOption is a function that configures the client
type ClientOption func(opts *clientOptions)

// WithAPIKey set the api key of the client
func WithAPIKey(apiKey string) ClientOption {
	return func(opts *clientOptions) {
		opts.apiKey = apiKey
	}
}

// WithBaseURL set the base url of the client
func WithBaseURL(baseURL string) ClientOption {
	return func(opts *clientOptions) {
		opts.baseURL = baseURL
	}
}

// WithDebug set the debug mode of the client
func WithDebug() ClientOption {
	return func(opts *clientOptions) {
		opts.debug = true
	}
}

// Client is the client for zhipu ai platform
type Client struct {
	client    *resty.Client
	keyID     string
	keySecret []byte
}

func (c *Client) createJWT() (token string, err error) {
	timestamp := time.Now().UnixMilli()
	exp := timestamp + time.Hour.Milliseconds()*24*7

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"api_key":   c.keyID,
		"timestamp": timestamp,
		"exp":       exp,
	})
	t.Header = map[string]interface{}{
		"alg":       "HS256",
		"sign_type": "SIGN",
	}

	return t.SignedString(c.keySecret)
}

// R creates a new resty request with the jwt token and context
func (c *Client) R(ctx context.Context) (req *resty.Request, err error) {
	var token string
	if token, err = c.createJWT(); err != nil {
		return
	}
	req = c.client.R().SetContext(ctx).SetHeader("Authorization", token)
	return
}

// NewClient creates a new client
// It will read the api key from the environment variable ZHIPUAI_API_KEY
// It will read the base url from the environment variable ZHIPUAI_BASE_URL
func NewClient(optFns ...ClientOption) (client *Client, err error) {
	var opts clientOptions
	for _, optFn := range optFns {
		optFn(&opts)
	}
	if opts.baseURL == "" {
		opts.baseURL = strings.TrimSpace(os.Getenv(envBaseURL))
	}
	if opts.baseURL == "" {
		opts.baseURL = defaultBaseURL
	}
	if opts.apiKey == "" {
		opts.apiKey = strings.TrimSpace(os.Getenv(envAPIKey))
	}
	if opts.apiKey == "" {
		err = ErrAPIKeyMissing
		return
	}
	ks := strings.SplitN(opts.apiKey, ".", 2)
	if len(ks) != 2 {
		err = ErrAPIKeyMalformed
		return
	}
	client = &Client{
		client:    resty.New().SetBaseURL(opts.baseURL),
		keyID:     ks[0],
		keySecret: []byte(ks[1]),
	}
	return
}
