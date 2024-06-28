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

type FineTuneError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type FineTuneItem struct {
	ID             string        `json:"id"`
	RequestID      string        `json:"request_id"`
	FineTunedModel string        `json:"fine_tuned_model"`
	Status         string        `json:"status"`
	Object         string        `json:"object"`
	TrainingFile   string        `json:"training_file"`
	ValidationFile string        `json:"validation_file"`
	Error          FineTuneError `json:"error"`
}

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

type FineTuneCreateResponse = FineTuneItem

func (c *Client) FineTuneCreateService(model string) *FineTuneCreateService {
	return &FineTuneCreateService{
		client: c,
		model:  model,
	}
}

func (s *FineTuneCreateService) SetModel(model string) *FineTuneCreateService {
	s.model = model
	return s
}

func (s *FineTuneCreateService) SetTrainingFile(trainingFile string) *FineTuneCreateService {
	s.trainingFile = trainingFile
	return s
}

func (s *FineTuneCreateService) SetValidationFile(validationFile string) *FineTuneCreateService {
	s.validationFile = &validationFile
	return s
}

func (s *FineTuneCreateService) SetLearningRateMultiplier(learningRateMultiplier float64) *FineTuneCreateService {
	s.learningRateMultiplier = &StringOr[float64]{}
	s.learningRateMultiplier.SetValue(learningRateMultiplier)
	return s
}

func (s *FineTuneCreateService) SetLearningRateMultiplierAuto() *FineTuneCreateService {
	s.learningRateMultiplier = &StringOr[float64]{}
	s.learningRateMultiplier.SetString(HyperParameterAuto)
	return s
}

func (s *FineTuneCreateService) SetBatchSize(batchSize int) *FineTuneCreateService {
	s.batchSize = &StringOr[int]{}
	s.batchSize.SetValue(batchSize)
	return s
}

func (s *FineTuneCreateService) SetBatchSizeAuto() *FineTuneCreateService {
	s.batchSize = &StringOr[int]{}
	s.batchSize.SetString(HyperParameterAuto)
	return s
}

func (s *FineTuneCreateService) SetNEpochs(nEpochs int) *FineTuneCreateService {
	s.nEpochs = &StringOr[int]{}
	s.nEpochs.SetValue(nEpochs)
	return s
}

func (s *FineTuneCreateService) SetNEpochsAuto() *FineTuneCreateService {
	s.nEpochs = &StringOr[int]{}
	s.nEpochs.SetString(HyperParameterAuto)
	return s
}

func (s *FineTuneCreateService) SetSuffix(suffix string) *FineTuneCreateService {
	s.suffix = &suffix
	return s
}

func (s *FineTuneCreateService) SetRequestID(requestID string) *FineTuneCreateService {
	s.requestID = &requestID
	return s
}

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

	if resp, err = s.client.request(ctx).SetBody(body).SetResult(&res).SetError(&apiError).Post("fine_tuning/jobs"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

type FineTuneEventListService struct {
	client *Client

	jobID string

	limit *int
	after *string
}

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

type FineTuneEventItem struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Object    string            `json:"object"`
	CreatedAt int64             `json:"created_at"`
	Data      FineTuneEventData `json:"data"`
}

type FineTuneEventListResponse struct {
	Data    []FineTuneEventItem `json:"data"`
	HasMore bool                `json:"has_more"`
	Object  string              `json:"object"`
}

func (c *Client) FineTuneEventListService(jobID string) *FineTuneEventListService {
	return &FineTuneEventListService{
		client: c,
		jobID:  jobID,
	}
}

func (s *FineTuneEventListService) SetJobID(jobID string) *FineTuneEventListService {
	s.jobID = jobID
	return s
}

func (s *FineTuneEventListService) SetLimit(limit int) *FineTuneEventListService {
	s.limit = &limit
	return s
}

func (s *FineTuneEventListService) SetAfter(after string) *FineTuneEventListService {
	s.after = &after
	return s
}

func (s *FineTuneEventListService) Do(ctx context.Context) (res FineTuneEventListResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	params := map[string]string{}
	if s.limit != nil {
		params["limit"] = strconv.Itoa(*s.limit)
	}
	if s.after != nil {
		params["after"] = *s.after
	}

	if resp, err = s.client.request(ctx).
		SetPathParam("job_id", s.jobID).
		SetQueryParams(params).
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

type FineTuneGetService struct {
	client *Client
	jobID  string
}

func (c *Client) FineTuneGetService(jobID string) *FineTuneGetService {
	return &FineTuneGetService{
		client: c,
		jobID:  jobID,
	}
}

func (s *FineTuneGetService) SetJobID(jobID string) *FineTuneGetService {
	s.jobID = jobID
	return s
}

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

type FineTuneListService struct {
	client *Client

	limit *int
	after *string
}

type FineTuneListResponse struct {
	Data   []FineTuneItem `json:"data"`
	Object string         `json:"object"`
}

func (c *Client) FineTuneListService() *FineTuneListService {
	return &FineTuneListService{
		client: c,
	}
}

func (s *FineTuneListService) SetLimit(limit int) *FineTuneListService {
	s.limit = &limit
	return s
}
func (s *FineTuneListService) SetAfter(after string) *FineTuneListService {
	s.after = &after
	return s
}

func (s *FineTuneListService) Do(ctx context.Context) (res FineTuneListResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	params := map[string]string{}
	if s.limit != nil {
		params["limit"] = strconv.Itoa(*s.limit)
	}
	if s.after != nil {
		params["after"] = *s.after
	}

	if resp, err = s.client.request(ctx).
		SetQueryParams(params).
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

type FineTuneDeleteService struct {
	client *Client
	jobID  string
}

func (c *Client) FineTuneDeleteService(jobID string) *FineTuneDeleteService {
	return &FineTuneDeleteService{
		client: c,
		jobID:  jobID,
	}
}

func (s *FineTuneDeleteService) SetJobID(jobID string) *FineTuneDeleteService {
	s.jobID = jobID
	return s
}

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

type FineTuneCancelService struct {
	client *Client
	jobID  string
}

func (c *Client) FineTuneCancelService(jobID string) *FineTuneCancelService {
	return &FineTuneCancelService{
		client: c,
		jobID:  jobID,
	}
}

func (s *FineTuneCancelService) SetJobID(jobID string) *FineTuneCancelService {
	s.jobID = jobID
	return s
}

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
