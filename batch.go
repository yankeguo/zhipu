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

	BatchCompletionWindow24h = "24h"
)

type BatchRequestCounts struct {
	Total     int64 `json:"total"`
	Completed int64 `json:"completed"`
	Failed    int64 `json:"failed"`
}

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

type BatchCreateService struct {
	client *Client

	inputFileID      string
	endpoint         string
	completionWindow string
	metadata         any
}

func (c *Client) BatchCreateService() *BatchCreateService {
	return &BatchCreateService{client: c}
}

func (s *BatchCreateService) SetInputFileID(inputFileID string) *BatchCreateService {
	s.inputFileID = inputFileID
	return s
}

func (s *BatchCreateService) SetEndpoint(endpoint string) *BatchCreateService {
	s.endpoint = endpoint
	return s
}

func (s *BatchCreateService) SetCompletionWindow(window string) *BatchCreateService {
	s.completionWindow = window
	return s
}

func (s *BatchCreateService) SetMetadata(metadata any) *BatchCreateService {
	s.metadata = metadata
	return s
}

func (s *BatchCreateService) Do(ctx context.Context) (res BatchItem, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).SetBody(M{
		"input_file_id":     s.inputFileID,
		"endpoint":          s.endpoint,
		"completion_window": s.completionWindow,
		"metadata":          s.metadata,
	}).SetResult(&res).SetError(&apiError).Post("batches"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
	}

	return
}

type BatchGetService struct {
	client  *Client
	batchID string
}

func (c *Client) BatchGetService(batchID string) *BatchGetService {
	return &BatchGetService{client: c, batchID: batchID}
}

func (s *BatchGetService) SetBatchID(batchID string) *BatchGetService {
	s.batchID = batchID
	return s
}

func (s *BatchGetService) Do(ctx context.Context) (res BatchItem, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetPathParam("batch_id", s.batchID).SetResult(&res).SetError(&apiError).
		Get("batches/{batch_id}"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
	}

	return
}

type BatchCancelService struct {
	client  *Client
	batchID string
}

func (c *Client) BatchCancelService(batchID string) *BatchCancelService {
	return &BatchCancelService{client: c, batchID: batchID}
}

func (s *BatchCancelService) SetBatchID(batchID string) *BatchCancelService {
	s.batchID = batchID
	return s
}

func (s *BatchCancelService) Do(ctx context.Context) (err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).SetBody(M{}).
		SetPathParam("batch_id", s.batchID).SetError(&apiError).
		Post("batches/{batch_id}/cancel"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
	}

	return
}

type BatchListService struct {
	client *Client

	after *string
	limit *int
}

type BatchListResponse struct {
	Object  string      `json:"object"`
	Data    []BatchItem `json:"data"`
	FirstID string      `json:"first_id"`
	LastID  string      `json:"last_id"`
	HasMore bool        `json:"has_more"`
}

func (c *Client) BatchListService() *BatchListService {
	return &BatchListService{client: c}
}

func (s *BatchListService) SetAfter(after string) *BatchListService {
	s.after = &after
	return s
}

func (s *BatchListService) SetLimit(limit int) *BatchListService {
	s.limit = &limit
	return s
}

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

	if resp, err = req.SetResult(&res).SetError(&apiError).Get("batches"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
	}

	return
}
