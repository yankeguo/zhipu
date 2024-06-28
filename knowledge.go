package zhipu

import (
	"context"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	KnowledgeEmbeddingIDEmbedding2 = 3
)

type KnowledgeCreateService struct {
	client      *Client
	embeddingID int
	name        string
	description *string
}

type KnowledgeCreateResponse = IDItem

func (c *Client) KnowledgeCreateService(name string) *KnowledgeCreateService {
	return &KnowledgeCreateService{
		client:      c,
		embeddingID: KnowledgeEmbeddingIDEmbedding2,
		name:        name,
	}
}

func (s *KnowledgeCreateService) SetEmbeddingID(embeddingID int) *KnowledgeCreateService {
	s.embeddingID = embeddingID
	return s
}

func (s *KnowledgeCreateService) SetName(name string) *KnowledgeCreateService {
	s.name = name
	return s
}

func (s *KnowledgeCreateService) SetDescription(description string) *KnowledgeCreateService {
	s.description = &description
	return s
}

func (s *KnowledgeCreateService) Do(ctx context.Context) (res KnowledgeCreateResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)
	body := M{
		"name":         s.name,
		"embedding_id": s.embeddingID,
	}
	if s.description != nil {
		body["description"] = *s.description
	}
	if resp, err = s.client.request(ctx).SetBody(body).SetResult(&res).SetError(&apiError).Post("knowledge"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

type KnowledgeEditService struct {
	client *Client

	knowledgeID string
	embeddingID *int
	name        *string
	description *string
}

func (c *Client) KnowledgeEditService(knowledgeID string) *KnowledgeEditService {
	return &KnowledgeEditService{
		client:      c,
		knowledgeID: knowledgeID,
	}
}

func (s *KnowledgeEditService) SetName(name string) *KnowledgeEditService {
	s.name = &name
	return s
}

func (s *KnowledgeEditService) SetEmbeddingID(embeddingID int) *KnowledgeEditService {
	s.embeddingID = &embeddingID
	return s
}

func (s *KnowledgeEditService) SetDescription(description string) *KnowledgeEditService {
	s.description = &description
	return s
}

func (s *KnowledgeEditService) Do(ctx context.Context) (err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)
	body := M{}
	if s.name != nil {
		body["name"] = *s.name
	}
	if s.description != nil {
		body["description"] = *s.description
	}
	if s.embeddingID != nil {
		body["embedding_id"] = *s.embeddingID
	}
	if resp, err = s.client.request(ctx).SetBody(body).SetPathParam("knowledge_id", s.knowledgeID).SetError(&apiError).Put("knowledge/{knowledge_id}"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

type KnowledgeListService struct {
	client *Client

	page *int
	size *int
}

type KnowledgeItem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Icon             string `json:"icon"`
	Background       string `json:"background"`
	EmbeddingID      int    `json:"embedding_id"`
	CustomIdentifier string `json:"custom_identifier"`
	WordNum          int64  `json:"word_num"`
	Length           int64  `json:"length"`
	DocumentSize     int64  `json:"document_size"`
}

type KnowledgeListResponse struct {
	List  []KnowledgeItem `json:"list"`
	Total int             `json:"total"`
}

func (c *Client) KnowledgeListService() *KnowledgeListService {
	return &KnowledgeListService{client: c}
}

func (s *KnowledgeListService) SetPage(page int) *KnowledgeListService {
	s.page = &page
	return s
}

func (s *KnowledgeListService) SetSize(size int) *KnowledgeListService {
	s.size = &size
	return s
}

func (s *KnowledgeListService) Do(ctx context.Context) (res KnowledgeListResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)
	req := s.client.request(ctx)
	if s.page != nil {
		req.SetQueryParam("page", strconv.Itoa(*s.page))
	}
	if s.size != nil {
		req.SetQueryParam("size", strconv.Itoa(*s.size))
	}
	if resp, err = req.SetResult(&res).SetError(&apiError).Get("knowledge"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

type KnowledgeDeleteService struct {
	client      *Client
	knowledgeID string
}

func (c *Client) KnowledgeDeleteService(knowledgeID string) *KnowledgeDeleteService {
	return &KnowledgeDeleteService{
		client:      c,
		knowledgeID: knowledgeID,
	}
}

func (s *KnowledgeDeleteService) SetKnowledgeID(knowledgeID string) *KnowledgeDeleteService {
	s.knowledgeID = knowledgeID
	return s
}

func (s *KnowledgeDeleteService) Do(ctx context.Context) (err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)
	if resp, err = s.client.request(ctx).SetPathParam("knowledge_id", s.knowledgeID).SetError(&apiError).Delete("knowledge/{knowledge_id}"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}
