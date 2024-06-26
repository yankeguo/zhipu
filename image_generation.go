package zhipu

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type ImageGenerationService struct {
	client *Client

	model  string
	prompt string
	userID string
}

type ImageGenerationResponse struct {
	Created int64      `json:"created"`
	Data    []ImageURL `json:"data"`
}

func (c *Client) ImageGenerationService(model string) *ImageGenerationService {
	return &ImageGenerationService{
		client: c,
		model:  model,
	}
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

func (s *ImageGenerationService) Do(ctx context.Context) (res ImageGenerationResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIError
	)

	body := M{
		"model":  s.model,
		"prompt": s.prompt,
	}

	if s.userID != "" {
		body["user_id"] = s.userID
	}

	if resp, err = s.client.R(ctx).SetBody(body).SetResult(&res).SetError(&apiError).Post("images/generations"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}
