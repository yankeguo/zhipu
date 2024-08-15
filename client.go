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

// BatchCreate creates a new BatchCreateService.
func (c *Client) BatchCreate() *BatchCreateService {
	return NewBatchCreateService(c)
}

// BatchGet creates a new BatchGetService.
func (c *Client) BatchGet(batchID string) *BatchGetService {
	return NewBatchGetService(c).SetBatchID(batchID)
}

// BatchCancel creates a new BatchCancelService.
func (c *Client) BatchCancel(batchID string) *BatchCancelService {
	return NewBatchCancelService(c).SetBatchID(batchID)
}

// BatchList creates a new BatchListService.
func (c *Client) BatchList() *BatchListService {
	return NewBatchListService(c)
}

// ChatCompletion creates a new ChatCompletionService.
func (c *Client) ChatCompletion(model string) *ChatCompletionService {
	return NewChatCompletionService(c).SetModel(model)
}

// Embedding embeds a list of text into a vector space.
func (c *Client) Embedding(model string) *EmbeddingService {
	return NewEmbeddingService(c).SetModel(model)
}

// FileCreate creates a new FileCreateService.
func (c *Client) FileCreate(purpose string) *FileCreateService {
	return NewFileCreateService(c).SetPurpose(purpose)
}

// FileEditService creates a new FileEditService.
func (c *Client) FileEdit(documentID string) *FileEditService {
	return NewFileEditService(c).SetDocumentID(documentID)
}

// FileList creates a new FileListService.
func (c *Client) FileList(purpose string) *FileListService {
	return NewFileListService(c).SetPurpose(purpose)
}

// FileDeleteService creates a new FileDeleteService.
func (c *Client) FileDelete(documentID string) *FileDeleteService {
	return NewFileDeleteService(c).SetDocumentID(documentID)
}

// FileGetService creates a new FileGetService.
func (c *Client) FileGet(documentID string) *FileGetService {
	return NewFileGetService(c).SetDocumentID(documentID)
}

// FileDownload creates a new FileDownloadService.
func (c *Client) FileDownload(fileID string) *FileDownloadService {
	return NewFileDownloadService(c).SetFileID(fileID)
}

// FineTuneCreate creates a new fine tune create service
func (c *Client) FineTuneCreate(model string) *FineTuneCreateService {
	return NewFineTuneCreateService(c).SetModel(model)
}

// FineTuneEventList creates a new fine tune event list service
func (c *Client) FineTuneEventList(jobID string) *FineTuneEventListService {
	return NewFineTuneEventListService(c).SetJobID(jobID)
}

// FineTuneGet creates a new fine tune get service
func (c *Client) FineTuneGet(jobID string) *FineTuneGetService {
	return NewFineTuneGetService(c).SetJobID(jobID)
}

// FineTuneList creates a new fine tune list service
func (c *Client) FineTuneList() *FineTuneListService {
	return NewFineTuneListService(c)
}

// FineTuneDelete creates a new fine tune delete service
func (c *Client) FineTuneDelete(jobID string) *FineTuneDeleteService {
	return NewFineTuneDeleteService(c).SetJobID(jobID)
}

// FineTuneCancel creates a new fine tune cancel service
func (c *Client) FineTuneCancel(jobID string) *FineTuneCancelService {
	return NewFineTuneCancelService(c).SetJobID(jobID)
}

// ImageGeneration creates a new image generation service
func (c *Client) ImageGeneration(model string) *ImageGenerationService {
	return NewImageGenerationService(c).SetModel(model)
}

// KnowledgeCreate creates a new knowledge create service
func (c *Client) KnowledgeCreate() *KnowledgeCreateService {
	return NewKnowledgeCreateService(c)
}

// KnowledgeEdit creates a new knowledge edit service
func (c *Client) KnowledgeEdit(knowledgeID string) *KnowledgeEditService {
	return NewKnowledgeEditService(c).SetKnowledgeID(knowledgeID)
}

// KnowledgeList list all the knowledge
func (c *Client) KnowledgeList() *KnowledgeListService {
	return NewKnowledgeListService(c)
}

// KnowledgeDelete creates a new knowledge delete service
func (c *Client) KnowledgeDelete(knowledgeID string) *KnowledgeDeleteService {
	return NewKnowledgeDeleteService(c).SetKnowledgeID(knowledgeID)
}

// KnowledgeGet creates a new knowledge get service
func (c *Client) KnowledgeCapacity() *KnowledgeCapacityService {
	return NewKnowledgeCapacityService(c)
}

// VideoGeneration creates a new video generation service
func (c *Client) VideoGeneration(model string) *VideoGenerationService {
	return NewVideoGenerationService(c).SetModel(model)
}

// AsyncResult creates a new async result get service
func (c *Client) AsyncResult(id string) *AsyncResultService {
	return NewAsyncResultService(c).SetID(id)
}
