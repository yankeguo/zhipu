package zhipu

import (
	"context"

	"github.com/go-resty/resty/v2"
)

const (
	VideoGenerationTaskStatusProcessing = "PROCESSING"
	VideoGenerationTaskStatusSuccess    = "SUCCESS"
	VideoGenerationTaskStatusFail       = "FAIL"
)

// VideoGenerationService creates a new video generation
type VideoGenerationService struct {
	client *Client

	model     string
	prompt    string
	userID    string
	imageURL  string
	requestID string
}

var (
	_ BatchSupport = &VideoGenerationService{}
)

// VideoGenerationResponse is the response of the VideoGenerationService
type VideoGenerationResponse struct {
	RequestID  string `json:"request_id"`
	ID         string `json:"id"`
	Model      string `json:"model"`
	TaskStatus string `json:"task_status"`
}

func NewVideoGenerationService(client *Client) *VideoGenerationService {
	return &VideoGenerationService{
		client: client,
	}
}

func (s *VideoGenerationService) BatchMethod() string {
	return "POST"
}

func (s *VideoGenerationService) BatchURL() string {
	return BatchEndpointV4VideosGenerations
}

func (s *VideoGenerationService) BatchBody() any {
	return s.buildBody()
}

// SetModel sets the model parameter
func (s *VideoGenerationService) SetModel(model string) *VideoGenerationService {
	s.model = model
	return s
}

// SetPrompt sets the prompt parameter
func (s *VideoGenerationService) SetPrompt(prompt string) *VideoGenerationService {
	s.prompt = prompt
	return s
}

// SetUserID sets the userID parameter
func (s *VideoGenerationService) SetUserID(userID string) *VideoGenerationService {
	s.userID = userID
	return s
}

// SetImageURL sets the imageURL parameter
func (s *VideoGenerationService) SetImageURL(imageURL string) *VideoGenerationService {
	s.imageURL = imageURL
	return s
}

// SetRequestID sets the requestID parameter
func (s *VideoGenerationService) SetRequestID(requestID string) *VideoGenerationService {
	s.requestID = requestID
	return s
}

func (s *VideoGenerationService) buildBody() M {
	body := M{
		"model":  s.model,
		"prompt": s.prompt,
	}
	if s.userID != "" {
		body["user_id"] = s.userID
	}
	if s.imageURL != "" {
		body["image_url"] = s.imageURL
	}
	if s.requestID != "" {
		body["request_id"] = s.requestID
	}
	return body
}

func (s *VideoGenerationService) Do(ctx context.Context) (res VideoGenerationResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	body := s.buildBody()

	if resp, err = s.client.request(ctx).
		SetBody(body).
		SetResult(&res).
		SetError(&apiError).
		Post("videos/generations"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}
