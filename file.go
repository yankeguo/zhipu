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

// FileCreateService is a service to create a file.
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

// FileCreateKnowledgeSuccessInfo is the success info of the FileCreateKnowledgeResponse.
type FileCreateKnowledgeSuccessInfo struct {
	Filename   string `json:"fileName"`
	DocumentID string `json:"documentId"`
}

// FileCreateKnowledgeFailedInfo is the failed info of the FileCreateKnowledgeResponse.
type FileCreateKnowledgeFailedInfo struct {
	Filename   string `json:"fileName"`
	FailReason string `json:"failReason"`
}

// FileCreateKnowledgeResponse is the response of the FileCreateService.
type FileCreateKnowledgeResponse struct {
	SuccessInfos []FileCreateKnowledgeSuccessInfo `json:"successInfos"`
	FailedInfos  []FileCreateKnowledgeFailedInfo  `json:"failedInfos"`
}

// FileCreateFineTuneResponse is the response of the FileCreateService.
type FileCreateFineTuneResponse struct {
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Object    string `json:"object"`
	Purpose   string `json:"purpose"`
	ID        string `json:"id"`
}

// FileCreateResponse is the response of the FileCreateService.
type FileCreateResponse struct {
	FileCreateFineTuneResponse
	FileCreateKnowledgeResponse
}

// NewFileCreateService creates a new FileCreateService.
func NewFileCreateService(client *Client) *FileCreateService {
	return &FileCreateService{client: client}
}

// SetLocalFile sets the local_file parameter of the FileCreateService.
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

	body := map[string]string{"purpose": s.purpose}

	if s.customSeparator != nil {
		body["custom_separator"] = *s.customSeparator
	}
	if s.sentenceSize != nil {
		body["sentence_size"] = strconv.Itoa(*s.sentenceSize)
	}
	if s.knowledgeID != nil {
		body["knowledge_id"] = *s.knowledgeID
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

	if file == nil {
		err = errors.New("no file specified")
		return
	}

	if resp, err = s.client.request(ctx).
		SetFileReader("file", filename, file).
		SetMultipartFormData(body).
		SetResult(&res).
		SetError(&apiError).
		Post("files"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}

// FileEditService is a service to edit a file.
type FileEditService struct {
	client *Client

	documentID string

	knowledgeType   *int
	customSeparator []string
	sentenceSize    *int
}

// NewFileEditService creates a new FileEditService.
func NewFileEditService(client *Client) *FileEditService {
	return &FileEditService{client: client}
}

// SetDocumentID sets the document_id parameter of the FileEditService.
func (s *FileEditService) SetDocumentID(documentID string) *FileEditService {
	s.documentID = documentID
	return s
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

// Do makes the request.
func (s *FileEditService) Do(ctx context.Context) (err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	body := M{}

	if s.knowledgeType != nil {
		body["knowledge_type"] = strconv.Itoa(*s.knowledgeType)
	}
	if len(s.customSeparator) > 0 {
		body["custom_separator"] = s.customSeparator
	}
	if s.sentenceSize != nil {
		body["sentence_size"] = strconv.Itoa(*s.sentenceSize)
	}

	if resp, err = s.client.request(ctx).
		SetPathParam("document_id", s.documentID).
		SetBody(body).
		SetError(&apiError).
		Put("document/{document_id}"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}

// FileListService is a service to list files.
type FileListService struct {
	client *Client

	purpose string

	knowledgeID *string
	page        *int
	limit       *int
	after       *string
	orderAsc    *bool
}

// FileFailInfo is the failed info of the FileListKnowledgeItem.
type FileFailInfo struct {
	EmbeddingCode int    `json:"embedding_code"`
	EmbeddingMsg  string `json:"embedding_msg"`
}

// FileListKnowledgeItem is the item of the FileListKnowledgeResponse.
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

// FileListKnowledgeResponse is the response of the FileListService.
type FileListKnowledgeResponse struct {
	Total int                     `json:"total"`
	List  []FileListKnowledgeItem `json:"list"`
}

// FileListFineTuneItem is the item of the FileListFineTuneResponse.
type FileListFineTuneItem struct {
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	ID        string `json:"id"`
	Object    string `json:"object"`
	Purpose   string `json:"purpose"`
}

// FileListFineTuneResponse is the response of the FileListService.
type FileListFineTuneResponse struct {
	Object string                 `json:"object"`
	Data   []FileListFineTuneItem `json:"data"`
}

// FileListResponse is the response of the FileListService.
type FileListResponse struct {
	FileListKnowledgeResponse
	FileListFineTuneResponse
}

// NewFileListService creates a new FileListService.
func NewFileListService(client *Client) *FileListService {
	return &FileListService{client: client}
}

// SetPurpose sets the purpose parameter of the FileListService.
func (s *FileListService) SetPurpose(purpose string) *FileListService {
	s.purpose = purpose
	return s
}

// SetKnowledgeID sets the knowledge_id parameter of the FileListService.
func (s *FileListService) SetKnowledgeID(knowledgeID string) *FileListService {
	s.knowledgeID = &knowledgeID
	return s
}

// SetPage sets the page parameter of the FileListService.
func (s *FileListService) SetPage(page int) *FileListService {
	s.page = &page
	return s
}

// SetLimit sets the limit parameter of the FileListService.
func (s *FileListService) SetLimit(limit int) *FileListService {
	s.limit = &limit
	return s
}

// SetAfter sets the after parameter of the FileListService.
func (s *FileListService) SetAfter(after string) *FileListService {
	s.after = &after
	return s
}

// SetOrder sets the order parameter of the FileListService.
func (s *FileListService) SetOrder(asc bool) *FileListService {
	s.orderAsc = &asc
	return s
}

// Do makes the request.
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
		SetQueryParams(m).
		SetResult(&res).
		SetError(&apiError).
		Get("files"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}

// FileDeleteService is a service to delete a file.
type FileDeleteService struct {
	client     *Client
	documentID string
}

// NewFileDeleteService creates a new FileDeleteService.
func NewFileDeleteService(client *Client) *FileDeleteService {
	return &FileDeleteService{client: client}
}

// SetDocumentID sets the document_id parameter of the FileDeleteService.
func (s *FileDeleteService) SetDocumentID(documentID string) *FileDeleteService {
	s.documentID = documentID
	return s
}

// Do makes the request.
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

// FileGetService is a service to get a file.
type FileGetService struct {
	client     *Client
	documentID string
}

// FileGetResponse is the response of the FileGetService.
type FileGetResponse = FileListKnowledgeItem

// NewFileGetService creates a new FileGetService.
func NewFileGetService(client *Client) *FileGetService {
	return &FileGetService{client: client}
}

// SetDocumentID sets the document_id parameter of the FileGetService.
func (s *FileGetService) SetDocumentID(documentID string) *FileGetService {
	s.documentID = documentID
	return s
}

// Do makes the request.
func (s *FileGetService) Do(ctx context.Context) (res FileGetResponse, err error) {
	var (
		resp     *resty.Response
		apiError APIErrorResponse
	)

	if resp, err = s.client.request(ctx).
		SetPathParam("document_id", s.documentID).
		SetResult(&res).
		SetError(&apiError).
		Get("document/{document_id}"); err != nil {
		return
	}

	if resp.IsError() {
		err = apiError
		return
	}

	return
}

// FileDownloadService is a service to download a file.
type FileDownloadService struct {
	client *Client

	fileID string

	writer   io.Writer
	filename string
}

// NewFileDownloadService creates a new FileDownloadService.
func NewFileDownloadService(client *Client) *FileDownloadService {
	return &FileDownloadService{client: client}
}

// SetFileID sets the file_id parameter of the FileDownloadService.
func (s *FileDownloadService) SetFileID(fileID string) *FileDownloadService {
	s.fileID = fileID
	return s
}

// SetOutput sets the output parameter of the FileDownloadService.
func (s *FileDownloadService) SetOutput(w io.Writer) *FileDownloadService {
	s.writer = w
	return s
}

// SetOutputFile sets the output_file parameter of the FileDownloadService.
func (s *FileDownloadService) SetOutputFile(filename string) *FileDownloadService {
	s.filename = filename
	return s
}

// Do makes the request.
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
