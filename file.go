package zhipu

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	FilePurposeFineTune  = "fine-tune"
	FilePurposeRetrieval = "retrieval"
	FilePurposeBatch     = "batch"

	KnowledgeTypeArticle                    = 1
	KnowledgeTypeQADocument                 = 2
	KnowledgeTypeQASpreadsheet              = 3
	KnowledgeTypeProductDatabaseSpreadsheet = 4
	KnowledgeTypeCustom                     = 5
)

type FileCreateService struct {
	client *Client

	purpose string

	localFile string
	file      io.Reader
	filename  string

	customSeparator *string
	sentenceSize    *int
	knowledgeID     *string
}

type FileCreateKnowledgeSuccessInfo struct {
	Filename   string `json:"fileName"`
	DocumentID string `json:"documentId"`
}

type FileCreateKnowledgeFailedInfo struct {
	Filename   string `json:"fileName"`
	FailReason string `json:"failReason"`
}

type FileCreateKnowledgeResponse struct {
	SuccessInfos []FileCreateKnowledgeSuccessInfo `json:"successInfos"`
	FailedInfos  []FileCreateKnowledgeFailedInfo  `json:"failedInfos"`
}

type FileCreateFineTuneResponse struct {
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Object    string `json:"object"`
	Purpose   string `json:"purpose"`
	ID        string `json:"id"`
}

type FileCreateResponse struct {
	FileCreateFineTuneResponse
	FileCreateKnowledgeResponse
}

// FileCreateService creates a new FileCreateService.
func (c *Client) FileCreateService(purpose string) *FileCreateService {
	return &FileCreateService{client: c, purpose: purpose}
}

func (s *FileCreateService) SetLocalFile(localFile string) *FileCreateService {
	s.localFile = localFile
	return s
}

// SetFile sets the file parameter of the FileCreateService.
func (s *FileCreateService) SetFile(file io.Reader, filename string) *FileCreateService {
	s.file = file
	s.filename = filename
	return s
}

// SetPurpose sets the purpose parameter of the FileCreateService.
func (s *FileCreateService) SetPurpose(purpose string) *FileCreateService {
	s.purpose = purpose
	return s
}

// SetCustomSeparator sets the custom_separator parameter of the FileCreateService.
func (s *FileCreateService) SetCustomSeparator(customSeparator string) *FileCreateService {
	s.customSeparator = &customSeparator
	return s
}

// SetSentenceSize sets the sentence_size parameter of the FileCreateService.
func (s *FileCreateService) SetSentenceSize(sentenceSize int) *FileCreateService {
	s.sentenceSize = &sentenceSize
	return s
}

// SetKnowledgeID sets the knowledge_id parameter of the FileCreateService.
func (s *FileCreateService) SetKnowledgeID(knowledgeID string) *FileCreateService {
	s.knowledgeID = &knowledgeID
	return s
}

// Do makes the request.
func (s *FileCreateService) Do(ctx context.Context) (res FileCreateResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	m := map[string]string{"purpose": s.purpose}

	if s.customSeparator != nil {
		m["custom_separator"] = *s.customSeparator
	}
	if s.sentenceSize != nil {
		m["sentence_size"] = strconv.Itoa(*s.sentenceSize)
	}
	if s.knowledgeID != nil {
		m["knowledge_id"] = *s.knowledgeID
	}

	file, filename := s.file, s.filename

	if file == nil && s.localFile != "" {
		var f *os.File
		if f, err = os.Open(s.localFile); err != nil {
			return
		}
		defer f.Close()

		file = f
		filename = filepath.Base(s.localFile)
	}

	if resp, err = s.client.request(ctx).
		SetResult(&res).
		SetError(&apiError).
		SetFileReader("file", filename, file).
		SetMultipartFormData(m).
		Post("files"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}

type FileEditService struct {
	client *Client

	documentID string

	knowledgeType   *int
	customSeparator []string
	sentenceSize    *int
}

// FileEditService creates a new FileEditService.
func (c *Client) FileEditService(documentID string) *FileEditService {
	return &FileEditService{client: c, documentID: documentID}
}

// SetKnowledgeType sets the knowledge_type parameter of the FileEditService.
func (s *FileEditService) SetKnowledgeType(knowledgeType int) *FileEditService {
	s.knowledgeType = &knowledgeType
	return s
}

// SetSentenceSize sets the sentence_size parameter of the FileEditService.
func (s *FileEditService) SetCustomSeparator(customSeparator ...string) *FileEditService {
	s.customSeparator = customSeparator
	return s
}

// SetSentenceSize sets the sentence_size parameter of the FileEditService.
func (s *FileEditService) SetSentenceSize(sentenceSize int) *FileEditService {
	s.sentenceSize = &sentenceSize
	return s
}

func (s *FileEditService) Do(ctx context.Context) (err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	m := M{}

	if s.knowledgeType != nil {
		m["knowledge_type"] = strconv.Itoa(*s.knowledgeType)
	}
	if len(s.customSeparator) > 0 {
		m["custom_separator"] = s.customSeparator
	}
	if s.sentenceSize != nil {
		m["sentence_size"] = strconv.Itoa(*s.sentenceSize)
	}

	if resp, err = s.client.request(ctx).
		SetPathParam("document_id", s.documentID).
		SetError(&apiError).
		SetBody(m).
		Put("document/{document_id}"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}

type FileListService struct {
	client *Client

	purpose string

	knowledgeID *string
	page        *int
	limit       *int
	after       *string
	orderAsc    *bool
}

type FileFailInfo struct {
	EmbeddingCode int    `json:"embedding_code"`
	EmbeddingMsg  string `json:"embedding_msg"`
}

type FileListKnowledgeItem struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	URL             string        `json:"url"`
	Length          int64         `json:"length"`
	SentenceSize    int64         `json:"sentence_size"`
	CustomSeparator []string      `json:"custom_separator"`
	EmbeddingStat   int           `json:"embedding_stat"`
	FailInfo        *FileFailInfo `json:"failInfo"`
	WordNum         int64         `json:"word_num"`
	ParseImage      int           `json:"parse_image"`
}

type FileListKnowledgeResponse struct {
	Total int                     `json:"total"`
	List  []FileListKnowledgeItem `json:"list"`
}

type FileListFineTuneItem struct {
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	ID        string `json:"id"`
	Object    string `json:"object"`
	Purpose   string `json:"purpose"`
}

type FileListFineTuneResponse struct {
	Object string                 `json:"object"`
	Data   []FileListFineTuneItem `json:"data"`
}

type FileListResponse struct {
	FileListKnowledgeResponse
	FileListFineTuneResponse
}

// FileListService creates a new FileListService.
func (c *Client) FileListService(purpose string) *FileListService {
	return &FileListService{client: c, purpose: purpose}
}

func (s *FileListService) SetPurpose(purpose string) *FileListService {
	s.purpose = purpose
	return s
}

func (s *FileListService) SetKnowledgeID(knowledgeID string) *FileListService {
	s.knowledgeID = &knowledgeID
	return s
}

func (s *FileListService) SetPage(page int) *FileListService {
	s.page = &page
	return s
}

func (s *FileListService) SetLimit(limit int) *FileListService {
	s.limit = &limit
	return s
}

func (s *FileListService) SetAfter(after string) *FileListService {
	s.after = &after
	return s
}

func (s *FileListService) SetOrder(asc bool) *FileListService {
	s.orderAsc = &asc
	return s
}

func (s *FileListService) Do(ctx context.Context) (res FileListResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	m := map[string]string{
		"purpose": s.purpose,
	}

	if s.knowledgeID != nil {
		m["knowledge_id"] = *s.knowledgeID
	}
	if s.page != nil {
		m["page"] = strconv.Itoa(*s.page)
	}
	if s.limit != nil {
		m["limit"] = strconv.Itoa(*s.limit)
	}
	if s.after != nil {
		m["after"] = *s.after
	}
	if s.orderAsc != nil {
		if *s.orderAsc {
			m["order"] = "asc"
		} else {
			m["order"] = "desc"
		}
	}

	if resp, err = s.client.request(ctx).
		SetResult(&res).
		SetError(&apiError).
		SetQueryParams(m).
		Get("files"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}

type FileDeleteService struct {
	client     *Client
	documentID string
}

// FileDeleteService creates a new FileDeleteService.
func (c *Client) FileDeleteService(documentID string) *FileDeleteService {
	return &FileDeleteService{client: c, documentID: documentID}
}

func (s *FileDeleteService) Do(ctx context.Context) (err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetPathParam("document_id", s.documentID).
		SetError(&apiError).
		Delete("document/{document_id}"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}

type FileGetService struct {
	client     *Client
	documentID string
}

type FileGetResponse = FileListKnowledgeItem

func (c *Client) FileGetService(documentID string) *FileGetService {
	return &FileGetService{client: c, documentID: documentID}
}

func (s *FileGetService) Do(ctx context.Context) (res FileGetResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetResult(&res).
		SetError(&apiError).
		SetPathParam("document_id", s.documentID).
		Get("document/{document_id}"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}

type FileDownloadService struct {
	client *Client

	fileID string

	writer   io.Writer
	filename string
}

func (c *Client) FileDownloadService(fileID string) *FileDownloadService {
	return &FileDownloadService{client: c, fileID: fileID}
}

func (s *FileDownloadService) SetOutput(w io.Writer) *FileDownloadService {
	s.writer = w
	return s
}

func (s *FileDownloadService) SetOutputFile(filename string) *FileDownloadService {
	s.filename = filename
	return s
}

func (s *FileDownloadService) Do(ctx context.Context) (err error) {
	var resp *resty.Response

	writer := s.writer

	if writer == nil && s.filename != "" {
		var f *os.File
		if f, err = os.Create(s.filename); err != nil {
			return
		}
		defer f.Close()

		writer = f
	}

	if writer == nil {
		return errors.New("no output specified")
	}

	if resp, err = s.client.request(ctx).
		SetDoNotParseResponse(true).
		SetPathParam("file_id", s.fileID).
		Get("files/{file_id}/content"); err != nil {
		return
	}
	defer resp.RawBody().Close()

	_, err = io.Copy(writer, resp.RawBody())

	return
}
