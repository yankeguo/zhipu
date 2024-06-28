package zhipu

import (
	"context"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	KnowledgeEmbeddingIDEmbedding2 = 3
)

// KnowledgeCreateService creates a new knowledge
type KnowledgeCreateService struct {
	client *Client

	embeddingID int
	name        string
	description *string
}

// KnowledgeCreateResponse is the response of the KnowledgeCreateService
type KnowledgeCreateResponse = IDItem

// NewKnowledgeCreateService creates a new KnowledgeCreateService
func NewKnowledgeCreateService(client *Client) *KnowledgeCreateService {
	return &KnowledgeCreateService{
		client: client,
	}
}

// SetEmbeddingID sets the embedding id of the knowledge
func (s *KnowledgeCreateService) SetEmbeddingID(embeddingID int) *KnowledgeCreateService {
	s.embeddingID = embeddingID
	return s
}

// SetName sets the name of the knowledge
func (s *KnowledgeCreateService) SetName(name string) *KnowledgeCreateService {
	s.name = name
	return s
}

// SetDescription sets the description of the knowledge
func (s *KnowledgeCreateService) SetDescription(description string) *KnowledgeCreateService {
	s.description = &description
	return s
}

// Do creates the knowledge
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
	if resp, err = s.client.request(ctx).
		SetBody(body).
		SetResult(&res).
		SetError(&apiError).
		Post("knowledge"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

// KnowledgeEditService edits a knowledge
type KnowledgeEditService struct {
	client *Client

	knowledgeID string

	embeddingID *int
	name        *string
	description *string
}

// NewKnowledgeEditService creates a new KnowledgeEditService
func NewKnowledgeEditService(client *Client) *KnowledgeEditService {
	return &KnowledgeEditService{
		client: client,
	}
}

// SetKnowledgeID sets the knowledge id
func (s *KnowledgeEditService) SetKnowledgeID(knowledgeID string) *KnowledgeEditService {
	s.knowledgeID = knowledgeID
	return s
}

// SetName sets the name of the knowledge
func (s *KnowledgeEditService) SetName(name string) *KnowledgeEditService {
	s.name = &name
	return s
}

// SetEmbeddingID sets the embedding id of the knowledge
func (s *KnowledgeEditService) SetEmbeddingID(embeddingID int) *KnowledgeEditService {
	s.embeddingID = &embeddingID
	return s
}

// SetDescription sets the description of the knowledge
func (s *KnowledgeEditService) SetDescription(description string) *KnowledgeEditService {
	s.description = &description
	return s
}

// Do edits the knowledge
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
	if resp, err = s.client.request(ctx).
		SetPathParam("knowledge_id", s.knowledgeID).
		SetBody(body).
		SetError(&apiError).
		Put("knowledge/{knowledge_id}"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

// KnowledgeListService lists the knowledge
type KnowledgeListService struct {
	client *Client

	page *int
	size *int
}

// KnowledgeItem is an item in the knowledge list
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

// KnowledgeListResponse is the response of the KnowledgeListService
type KnowledgeListResponse struct {
	List  []KnowledgeItem `json:"list"`
	Total int             `json:"total"`
}

// NewKnowledgeListService creates a new KnowledgeListService
func NewKnowledgeListService(client *Client) *KnowledgeListService {
	return &KnowledgeListService{client: client}
}

// SetPage sets the page of the knowledge list
func (s *KnowledgeListService) SetPage(page int) *KnowledgeListService {
	s.page = &page
	return s
}

// SetSize sets the size of the knowledge list
func (s *KnowledgeListService) SetSize(size int) *KnowledgeListService {
	s.size = &size
	return s
}

// Do lists the knowledge
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
	if resp, err = req.
		SetResult(&res).
		SetError(&apiError).
		Get("knowledge"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

// KnowledgeDeleteService deletes a knowledge
type KnowledgeDeleteService struct {
	client *Client

	knowledgeID string
}

// NewKnowledgeDeleteService creates a new KnowledgeDeleteService
func NewKnowledgeDeleteService(client *Client) *KnowledgeDeleteService {
	return &KnowledgeDeleteService{
		client: client,
	}
}

// SetKnowledgeID sets the knowledge id
func (s *KnowledgeDeleteService) SetKnowledgeID(knowledgeID string) *KnowledgeDeleteService {
	s.knowledgeID = knowledgeID
	return s
}

// Do deletes the knowledge
func (s *KnowledgeDeleteService) Do(ctx context.Context) (err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)
	if resp, err = s.client.request(ctx).
		SetPathParam("knowledge_id", s.knowledgeID).
		SetError(&apiError).
		Delete("knowledge/{knowledge_id}"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}

// KnowledgeCapacityService query the capacity of the knowledge
type KnowledgeCapacityService struct {
	client *Client
}

// KnowledgeCapacityItem is an item in the knowledge capacity
type KnowledgeCapacityItem struct {
	WordNum int64 `json:"word_num"`
	Length  int64 `json:"length"`
}

// KnowledgeCapacityResponse is the response of the KnowledgeCapacityService
type KnowledgeCapacityResponse struct {
	Used  KnowledgeCapacityItem `json:"used"`
	Total KnowledgeCapacityItem `json:"total"`
}

// SetKnowledgeID sets the knowledge id
func NewKnowledgeCapacityService(client *Client) *KnowledgeCapacityService {
	return &KnowledgeCapacityService{client: client}
}

// Do query the capacity of the knowledge
func (s *KnowledgeCapacityService) Do(ctx context.Context) (res KnowledgeCapacityResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)
	if resp, err = s.client.request(ctx).
		SetResult(&res).
		SetError(&apiError).
		Get("knowledge/capacity"); err != nil {
		return
	}
	if resp.IsError() {
		err = apiError
		return
	}
	return
}
