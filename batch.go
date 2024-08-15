package zhipu

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	BatchEndpointV4ChatCompletions   = "/v4/chat/completions"
	BatchEndpointV4ImagesGenerations = "/v4/images/generations"
	BatchEndpointV4Embeddings        = "/v4/embeddings"
	BatchEndpointV4VideosGenerations = "/v4/videos/generations"

	BatchCompletionWindow24h = "24h"
)

// BatchRequestCounts represents the counts of the batch requests.
type BatchRequestCounts struct {
	Total     int64 `json:"total"`
	Completed int64 `json:"completed"`
	Failed    int64 `json:"failed"`
}

// BatchItem represents a batch item.
type BatchItem struct {
	ID               string             `json:"id"`
	Object           any                `json:"object"`
	Endpoint         string             `json:"endpoint"`
	InputFileID      string             `json:"input_file_id"`
	CompletionWindow string             `json:"completion_window"`
	Status           string             `json:"status"`
	OutputFileID     string             `json:"output_file_id"`
	ErrorFileID      string             `json:"error_file_id"`
	CreatedAt        int64              `json:"created_at"`
	InProgressAt     int64              `json:"in_progress_at"`
	ExpiresAt        int64              `json:"expires_at"`
	FinalizingAt     int64              `json:"finalizing_at"`
	CompletedAt      int64              `json:"completed_at"`
	FailedAt         int64              `json:"failed_at"`
	ExpiredAt        int64              `json:"expired_at"`
	CancellingAt     int64              `json:"cancelling_at"`
	CancelledAt      int64              `json:"cancelled_at"`
	RequestCounts    BatchRequestCounts `json:"request_counts"`
	Metadata         json.RawMessage    `json:"metadata"`
}

// BatchCreateService is a service to create a batch.
type BatchCreateService struct {
	client *Client

	inputFileID      string
	endpoint         string
	completionWindow string
	metadata         any
}

// NewBatchCreateService creates a new BatchCreateService.
func NewBatchCreateService(client *Client) *BatchCreateService {
	return &BatchCreateService{client: client}
}

// SetInputFileID sets the input file id for the batch.
func (s *BatchCreateService) SetInputFileID(inputFileID string) *BatchCreateService {
	s.inputFileID = inputFileID
	return s
}

// SetEndpoint sets the endpoint for the batch.
func (s *BatchCreateService) SetEndpoint(endpoint string) *BatchCreateService {
	s.endpoint = endpoint
	return s
}

// SetCompletionWindow sets the completion window for the batch.
func (s *BatchCreateService) SetCompletionWindow(window string) *BatchCreateService {
	s.completionWindow = window
	return s
}

// SetMetadata sets the metadata for the batch.
func (s *BatchCreateService) SetMetadata(metadata any) *BatchCreateService {
	s.metadata = metadata
	return s
}

// Do executes the batch create service.
func (s *BatchCreateService) Do(ctx context.Context) (res BatchItem, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetBody(M{
			"input_file_id":     s.inputFileID,
			"endpoint":          s.endpoint,
			"completion_window": s.completionWindow,
			"metadata":          s.metadata,
		}).
		SetResult(&res).
		SetError(&apiError).
		Post("batches"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
	}

	return
}

// BatchGetService is a service to get a batch.
type BatchGetService struct {
	client  *Client
	batchID string
}

// BatchGetResponse represents the response of the batch get service.
type BatchGetResponse = BatchItem

// NewBatchGetService creates a new BatchGetService.
func NewBatchGetService(client *Client) *BatchGetService {
	return &BatchGetService{client: client}
}

// SetBatchID sets the batch id for the batch get service.
func (s *BatchGetService) SetBatchID(batchID string) *BatchGetService {
	s.batchID = batchID
	return s
}

// Do executes the batch get service.
func (s *BatchGetService) Do(ctx context.Context) (res BatchGetResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetPathParam("batch_id", s.batchID).
		SetResult(&res).
		SetError(&apiError).
		Get("batches/{batch_id}"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
	}

	return
}

// BatchCancelService is a service to cancel a batch.
type BatchCancelService struct {
	client  *Client
	batchID string
}

// NewBatchCancelService creates a new BatchCancelService.
func NewBatchCancelService(client *Client) *BatchCancelService {
	return &BatchCancelService{client: client}
}

// SetBatchID sets the batch id for the batch cancel service.
func (s *BatchCancelService) SetBatchID(batchID string) *BatchCancelService {
	s.batchID = batchID
	return s
}

// Do executes the batch cancel service.
func (s *BatchCancelService) Do(ctx context.Context) (err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetPathParam("batch_id", s.batchID).
		SetBody(M{}).
		SetError(&apiError).
		Post("batches/{batch_id}/cancel"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
	}

	return
}

// BatchListService is a service to list batches.
type BatchListService struct {
	client *Client

	after *string
	limit *int
}

// BatchListResponse represents the response of the batch list service.
type BatchListResponse struct {
	Object  string      `json:"object"`
	Data    []BatchItem `json:"data"`
	FirstID string      `json:"first_id"`
	LastID  string      `json:"last_id"`
	HasMore bool        `json:"has_more"`
}

// NewBatchListService creates a new BatchListService.
func NewBatchListService(client *Client) *BatchListService {
	return &BatchListService{client: client}
}

// SetAfter sets the after cursor for the batch list service.
func (s *BatchListService) SetAfter(after string) *BatchListService {
	s.after = &after
	return s
}

// SetLimit sets the limit for the batch list service.
func (s *BatchListService) SetLimit(limit int) *BatchListService {
	s.limit = &limit
	return s
}

// Do executes the batch list service.
func (s *BatchListService) Do(ctx context.Context) (res BatchListResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	req := s.client.request(ctx)
	if s.after != nil {
		req.SetQueryParam("after", *s.after)
	}
	if s.limit != nil {
		req.SetQueryParam("limit", strconv.Itoa(*s.limit))
	}

	if resp, err = req.
		SetResult(&res).
		SetError(&apiError).
		Get("batches"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
	}

	return
}
