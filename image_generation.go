package zhipu

import (
	"context"

	"github.com/go-resty/resty/v2"
)

// ImageGenerationService creates a new image generation
type ImageGenerationService struct {
	client *Client

	model  string
	prompt string
	userID string
}

var (
	_ BatchSupport = &ImageGenerationService{}
)

// ImageGenerationResponse is the response of the ImageGenerationService
type ImageGenerationResponse struct {
	Created int64     `json:"created"`
	Data    []URLItem `json:"data"`
}

// NewImageGenerationService creates a new ImageGenerationService
func NewImageGenerationService(client *Client) *ImageGenerationService {
	return &ImageGenerationService{
		client: client,
	}
}

func (s *ImageGenerationService) BatchMethod() string {
	return "POST"
}

func (s *ImageGenerationService) BatchURL() string {
	return BatchEndpointV4ImagesGenerations
}

func (s *ImageGenerationService) BatchBody() any {
	return s.buildBody()
}

// SetModel sets the model parameter
func (s *ImageGenerationService) SetModel(model string) *ImageGenerationService {
	s.model = model
	return s
}

// SetPrompt sets the prompt parameter
func (s *ImageGenerationService) SetPrompt(prompt string) *ImageGenerationService {
	s.prompt = prompt
	return s
}

// SetUserID sets the userID parameter
func (s *ImageGenerationService) SetUserID(userID string) *ImageGenerationService {
	s.userID = userID
	return s
}

func (s *ImageGenerationService) buildBody() M {
	body := M{
		"model":  s.model,
		"prompt": s.prompt,
	}

	if s.userID != "" {
		body["user_id"] = s.userID
	}

	return body
}

func (s *ImageGenerationService) Do(ctx context.Context) (res ImageGenerationResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	body := s.buildBody()

	if resp, err = s.client.request(ctx).
		SetBody(body).
		SetResult(&res).
		SetError(&apiError).
		Post("images/generations"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}
