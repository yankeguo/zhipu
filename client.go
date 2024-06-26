package zhipu

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	envAPIKey  = "ZHIPUAI_API_KEY"
	envBaseURL = "ZHIPUAI_BASE_URL"
	envDebug   = "ZHIPUAI_DEBUG"

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
	client  *http.Client
	debug   *bool
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

// WithHTTPClient set the http client of the client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(opts *clientOptions) {
		opts.client = client
	}
}

// WithDebug set the debug mode of the client
func WithDebug(debug bool) ClientOption {
	return func(opts *clientOptions) {
		opts.debug = new(bool)
		*opts.debug = debug
	}
}

// Client is the client for zhipu ai platform
type Client struct {
	client    *resty.Client
	debug     bool
	keyID     string
	keySecret []byte
}

func (c *Client) createJWT() string {
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

	token, err := t.SignedString(c.keySecret)
	if err != nil {
		panic(err)
	}
	return token
}

// request creates a new resty request with the jwt token and context
func (c *Client) request(ctx context.Context) *resty.Request {
	return c.client.R().SetContext(ctx).SetHeader("Authorization", c.createJWT())
}

// NewClient creates a new client
// It will read the api key from the environment variable ZHIPUAI_API_KEY
// It will read the base url from the environment variable ZHIPUAI_BASE_URL
func NewClient(optFns ...ClientOption) (client *Client, err error) {
	var opts clientOptions
	for _, optFn := range optFns {
		optFn(&opts)
	}
	// base url
	if opts.baseURL == "" {
		opts.baseURL = strings.TrimSpace(os.Getenv(envBaseURL))
	}
	if opts.baseURL == "" {
		opts.baseURL = defaultBaseURL
	}
	// api key
	if opts.apiKey == "" {
		opts.apiKey = strings.TrimSpace(os.Getenv(envAPIKey))
	}
	if opts.apiKey == "" {
		err = ErrAPIKeyMissing
		return
	}
	// debug
	if opts.debug == nil {
		if debugStr := strings.TrimSpace(os.Getenv(envDebug)); debugStr != "" {
			if debug, err1 := strconv.ParseBool(debugStr); err1 == nil {
				opts.debug = &debug
			}
		}
	}

	keyComponents := strings.SplitN(opts.apiKey, ".", 2)

	if len(keyComponents) != 2 {
		err = ErrAPIKeyMalformed
		return
	}

	client = &Client{
		keyID:     keyComponents[0],
		keySecret: []byte(keyComponents[1]),
	}

	if opts.client == nil {
		client.client = resty.New()
	} else {
		client.client = resty.NewWithClient(opts.client)
	}

	client.client = client.client.SetBaseURL(opts.baseURL)

	if opts.debug != nil {
		client.client.SetDebug(*opts.debug)
		client.debug = *opts.debug
	}
	return
}
