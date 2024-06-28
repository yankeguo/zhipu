package zhipu

import (
	"context"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	HyperParameterAuto = "auto"

	FineTuneStatusCreate          = "create"
	FineTuneStatusValidatingFiles = "validating_files"
	FineTuneStatusQueued          = "queued"
	FineTuneStatusRunning         = "running"
	FineTuneStatusSucceeded       = "succeeded"
	FineTuneStatusFailed          = "failed"
	FineTuneStatusCancelled       = "cancelled"
)

// FineTuneItem is the item of the FineTune
type FineTuneItem struct {
	ID             string   `json:"id"`
	RequestID      string   `json:"request_id"`
	FineTunedModel string   `json:"fine_tuned_model"`
	Status         string   `json:"status"`
	Object         string   `json:"object"`
	TrainingFile   string   `json:"training_file"`
	ValidationFile string   `json:"validation_file"`
	Error          APIError `json:"error"`
}

// FineTuneCreateService creates a new fine tune
type FineTuneCreateService struct {
	client *Client

	model          string
	trainingFile   string
	validationFile *string

	learningRateMultiplier *StringOr[float64]
	batchSize              *StringOr[int]
	nEpochs                *StringOr[int]

	suffix    *string
	requestID *string
}

// FineTuneCreateResponse is the response of the FineTuneCreateService
type FineTuneCreateResponse = FineTuneItem

// NewFineTuneCreateService creates a new FineTuneCreateService
func NewFineTuneCreateService(client *Client) *FineTuneCreateService {
	return &FineTuneCreateService{
		client: client,
	}
}

// SetModel sets the model parameter
func (s *FineTuneCreateService) SetModel(model string) *FineTuneCreateService {
	s.model = model
	return s
}

// SetTrainingFile sets the trainingFile parameter
func (s *FineTuneCreateService) SetTrainingFile(trainingFile string) *FineTuneCreateService {
	s.trainingFile = trainingFile
	return s
}

// SetValidationFile sets the validationFile parameter
func (s *FineTuneCreateService) SetValidationFile(validationFile string) *FineTuneCreateService {
	s.validationFile = &validationFile
	return s
}

// SetLearningRateMultiplier sets the learningRateMultiplier parameter
func (s *FineTuneCreateService) SetLearningRateMultiplier(learningRateMultiplier float64) *FineTuneCreateService {
	s.learningRateMultiplier = &StringOr[float64]{}
	s.learningRateMultiplier.SetValue(learningRateMultiplier)
	return s
}

// SetLearningRateMultiplierAuto sets the learningRateMultiplier parameter to auto
func (s *FineTuneCreateService) SetLearningRateMultiplierAuto() *FineTuneCreateService {
	s.learningRateMultiplier = &StringOr[float64]{}
	s.learningRateMultiplier.SetString(HyperParameterAuto)
	return s
}

// SetBatchSize sets the batchSize parameter
func (s *FineTuneCreateService) SetBatchSize(batchSize int) *FineTuneCreateService {
	s.batchSize = &StringOr[int]{}
	s.batchSize.SetValue(batchSize)
	return s
}

// SetBatchSizeAuto sets the batchSize parameter to auto
func (s *FineTuneCreateService) SetBatchSizeAuto() *FineTuneCreateService {
	s.batchSize = &StringOr[int]{}
	s.batchSize.SetString(HyperParameterAuto)
	return s
}

// SetNEpochs sets the nEpochs parameter
func (s *FineTuneCreateService) SetNEpochs(nEpochs int) *FineTuneCreateService {
	s.nEpochs = &StringOr[int]{}
	s.nEpochs.SetValue(nEpochs)
	return s
}

// SetNEpochsAuto sets the nEpochs parameter to auto
func (s *FineTuneCreateService) SetNEpochsAuto() *FineTuneCreateService {
	s.nEpochs = &StringOr[int]{}
	s.nEpochs.SetString(HyperParameterAuto)
	return s
}

// SetSuffix sets the suffix parameter
func (s *FineTuneCreateService) SetSuffix(suffix string) *FineTuneCreateService {
	s.suffix = &suffix
	return s
}

// SetRequestID sets the requestID parameter
func (s *FineTuneCreateService) SetRequestID(requestID string) *FineTuneCreateService {
	s.requestID = &requestID
	return s
}

// Do makes the request
func (s *FineTuneCreateService) Do(ctx context.Context) (res FineTuneCreateResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	body := M{
		"model":         s.model,
		"training_file": s.trainingFile,
	}

	if s.validationFile != nil {
		body["validation_file"] = *s.validationFile
	}
	if s.suffix != nil {
		body["suffix"] = *s.suffix
	}
	if s.requestID != nil {
		body["request_id"] = *s.requestID
	}
	if s.learningRateMultiplier != nil || s.batchSize != nil || s.nEpochs != nil {
		hp := M{}
		if s.learningRateMultiplier != nil {
			hp["learning_rate_multiplier"] = s.learningRateMultiplier
		}
		if s.batchSize != nil {
			hp["batch_size"] = s.batchSize
		}
		if s.nEpochs != nil {
			hp["n_epochs"] = s.nEpochs
		}
		body["hyperparameters"] = hp
	}

	if resp, err = s.client.request(ctx).
		SetBody(body).
		SetResult(&res).
		SetError(&apiError).
		Post("fine_tuning/jobs"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

// FineTuneEventListService creates a new fine tune event list
type FineTuneEventListService struct {
	client *Client

	jobID string

	limit *int
	after *string
}

// FineTuneEventData is the data of the FineTuneEventItem
type FineTuneEventData struct {
	Acc           float64 `json:"acc"`
	Loss          float64 `json:"loss"`
	CurrentSteps  int64   `json:"current_steps"`
	RemainingTime string  `json:"remaining_time"`
	ElapsedTime   string  `json:"elapsed_time"`
	TotalSteps    int64   `json:"total_steps"`
	Epoch         int64   `json:"epoch"`
	TrainedTokens int64   `json:"trained_tokens"`
	LearningRate  float64 `json:"learning_rate"`
}

// FineTuneEventItem is the item of the FineTuneEventListResponse
type FineTuneEventItem struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Object    string            `json:"object"`
	CreatedAt int64             `json:"created_at"`
	Data      FineTuneEventData `json:"data"`
}

// FineTuneEventListResponse is the response of the FineTuneEventListService
type FineTuneEventListResponse struct {
	Data    []FineTuneEventItem `json:"data"`
	HasMore bool                `json:"has_more"`
	Object  string              `json:"object"`
}

// NewFineTuneEventListService creates a new FineTuneEventListService
func NewFineTuneEventListService(client *Client) *FineTuneEventListService {
	return &FineTuneEventListService{
		client: client,
	}
}

// SetJobID sets the jobID parameter
func (s *FineTuneEventListService) SetJobID(jobID string) *FineTuneEventListService {
	s.jobID = jobID
	return s
}

// SetLimit sets the limit parameter
func (s *FineTuneEventListService) SetLimit(limit int) *FineTuneEventListService {
	s.limit = &limit
	return s
}

// SetAfter sets the after parameter
func (s *FineTuneEventListService) SetAfter(after string) *FineTuneEventListService {
	s.after = &after
	return s
}

// Do makes the request
func (s *FineTuneEventListService) Do(ctx context.Context) (res FineTuneEventListResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	req := s.client.request(ctx)

	if s.limit != nil {
		req.SetQueryParam("limit", strconv.Itoa(*s.limit))
	}
	if s.after != nil {
		req.SetQueryParam("after", *s.after)
	}

	if resp, err = req.
		SetPathParam("job_id", s.jobID).
		SetResult(&res).
		SetError(&apiError).
		Get("fine_tuning/jobs/{job_id}/events"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

// FineTuneGetService creates a new fine tune get
type FineTuneGetService struct {
	client *Client
	jobID  string
}

// NewFineTuneGetService creates a new FineTuneGetService
func NewFineTuneGetService(client *Client) *FineTuneGetService {
	return &FineTuneGetService{
		client: client,
	}
}

// SetJobID sets the jobID parameter
func (s *FineTuneGetService) SetJobID(jobID string) *FineTuneGetService {
	s.jobID = jobID
	return s
}

// Do makes the request
func (s *FineTuneGetService) Do(ctx context.Context) (res FineTuneItem, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetPathParam("job_id", s.jobID).
		SetResult(&res).
		SetError(&apiError).
		Get("fine_tuning/jobs/{job_id}"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

// FineTuneListService creates a new fine tune list
type FineTuneListService struct {
	client *Client

	limit *int
	after *string
}

// FineTuneListResponse is the response of the FineTuneListService
type FineTuneListResponse struct {
	Data   []FineTuneItem `json:"data"`
	Object string         `json:"object"`
}

// NewFineTuneListService creates a new FineTuneListService
func NewFineTuneListService(client *Client) *FineTuneListService {
	return &FineTuneListService{
		client: client,
	}
}

// SetLimit sets the limit parameter
func (s *FineTuneListService) SetLimit(limit int) *FineTuneListService {
	s.limit = &limit
	return s
}

// SetAfter sets the after parameter
func (s *FineTuneListService) SetAfter(after string) *FineTuneListService {
	s.after = &after
	return s
}

// Do makes the request
func (s *FineTuneListService) Do(ctx context.Context) (res FineTuneListResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	req := s.client.request(ctx)
	if s.limit != nil {
		req.SetQueryParam("limit", strconv.Itoa(*s.limit))
	}
	if s.after != nil {
		req.SetQueryParam("after", *s.after)
	}

	if resp, err = req.
		SetResult(&res).
		SetError(&apiError).
		Get("fine_tuning/jobs"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

// FineTuneDeleteService creates a new fine tune delete
type FineTuneDeleteService struct {
	client *Client
	jobID  string
}

// NewFineTuneDeleteService creates a new FineTuneDeleteService
func NewFineTuneDeleteService(client *Client) *FineTuneDeleteService {
	return &FineTuneDeleteService{
		client: client,
	}
}

// SetJobID sets the jobID parameter
func (s *FineTuneDeleteService) SetJobID(jobID string) *FineTuneDeleteService {
	s.jobID = jobID
	return s
}

// Do makes the request
func (s *FineTuneDeleteService) Do(ctx context.Context) (res FineTuneItem, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetPathParam("job_id", s.jobID).
		SetResult(&res).
		SetError(&apiError).
		Delete("fine_tuning/jobs/{job_id}"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

// FineTuneCancelService creates a new fine tune cancel
type FineTuneCancelService struct {
	client *Client
	jobID  string
}

// NewFineTuneCancelService creates a new FineTuneCancelService
func NewFineTuneCancelService(client *Client) *FineTuneCancelService {
	return &FineTuneCancelService{
		client: client,
	}
}

// SetJobID sets the jobID parameter
func (s *FineTuneCancelService) SetJobID(jobID string) *FineTuneCancelService {
	s.jobID = jobID
	return s
}

// Do makes the request
func (s *FineTuneCancelService) Do(ctx context.Context) (res FineTuneItem, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetPathParam("job_id", s.jobID).
		SetResult(&res).
		SetError(&apiError).
		Post("fine_tuning/jobs/{job_id}/cancel"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}
