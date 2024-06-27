package zhipu

import (
	"context"

	"github.com/go-resty/resty/v2"
)

// EmbeddingData is the data for each embedding.
type EmbeddingData struct {
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
	Object    string    `json:"object"`
}

// EmbeddingResponse is the response from the embedding service.
type EmbeddingResponse struct {
	Model  string              `json:"model"`
	Data   []EmbeddingData     `json:"data"`
	Object string              `json:"object"`
	Usage  ChatCompletionUsage `json:"usage"`
}

// EmbeddingService embeds a list of text into a vector space.
type EmbeddingService struct {
	client *Client

	model string
	input string
}

// EmbeddingService embeds a list of text into a vector space.
func (c *Client) EmbeddingService(model string) *EmbeddingService {
	return &EmbeddingService{client: c, model: model}
}

// SetModel sets the model to use for the embedding.
func (s *EmbeddingService) SetModel(model string) *EmbeddingService {
	s.model = model
	return s
}

// SetInput sets the input text to embed.
func (s *EmbeddingService) SetInput(input string) *EmbeddingService {
	s.input = input
	return s
}

func (s *EmbeddingService) Do(ctx context.Context) (res EmbeddingResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).SetBody(M{"model": s.model, "input": s.input}).
		SetResult(&res).
		SetError(&apiError).
		Post("embeddings"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}
