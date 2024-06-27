package zhipu

import (
	"context"
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
