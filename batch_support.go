package zhipu

import (
	"encoding/json"
	"io"
)

// BatchSupport is the interface for services with batch support.
type BatchSupport interface {
	BatchMethod() string
	BatchURL() string
	BatchBody() any
}

// BatchFileWriter is a writer for batch files.
type BatchFileWriter struct {
	w  io.Writer
	je *json.Encoder
}

// NewBatchFileWriter creates a new BatchFileWriter.
func NewBatchFileWriter(w io.Writer) *BatchFileWriter {
	return &BatchFileWriter{w: w, je: json.NewEncoder(w)}
}

// Write writes a batch file.
func (b *BatchFileWriter) Write(customID string, s BatchSupport) error {
	return b.je.Encode(M{
		"custom_id": customID,
		"method":    s.BatchMethod(),
		"url":       s.BatchURL(),
		"body":      s.BatchBody(),
	})
}

// BatchResultResponse is the response of a batch result.
type BatchResultResponse[T any] struct {
	StatusCode int `json:"status_code"`
	Body       T   `json:"body"`
}

// BatchResult is the result of a batch.
type BatchResult[T any] struct {
	ID       string                 `json:"id"`
	CustomID string                 `json:"custom_id"`
	Response BatchResultResponse[T] `json:"response"`
}

// BatchResultReader reads batch results.
type BatchResultReader[T any] struct {
	r  io.Reader
	jd *json.Decoder
}

// NewBatchResultReader creates a new BatchResultReader.
func NewBatchResultReader[T any](r io.Reader) *BatchResultReader[T] {
	return &BatchResultReader[T]{r: r, jd: json.NewDecoder(r)}
}

// Read reads a batch result.
func (r *BatchResultReader[T]) Read(out *BatchResult[T]) error {
	return r.jd.Decode(out)
}
