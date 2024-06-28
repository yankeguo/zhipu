package zhipu

import (
	"encoding/json"
	"io"
)

type BatchSupport interface {
	BatchMethod() string
	BatchURL() string
	BatchBody() any
}

type BatchFileWriter struct {
	w  io.Writer
	je *json.Encoder
}

func NewBatchFileWriter(w io.Writer) *BatchFileWriter {
	return &BatchFileWriter{w: w, je: json.NewEncoder(w)}
}

func (b *BatchFileWriter) Write(customID string, s BatchSupport) error {
	return b.je.Encode(M{
		"custom_id": customID,
		"method":    s.BatchMethod(),
		"url":       s.BatchURL(),
		"body":      s.BatchBody(),
	})
}

type BatchResultResponse[T any] struct {
	StatusCode int `json:"status_code"`
	Body       T   `json:"body"`
}

type BatchResult[T any] struct {
	ID       string                 `json:"id"`
	CustomID string                 `json:"custom_id"`
	Response BatchResultResponse[T] `json:"response"`
}

type BatchResultReader[T any] struct {
	r  io.Reader
	jd *json.Decoder
}

func NewBatchResultReader[T any](r io.Reader) *BatchResultReader[T] {
	return &BatchResultReader[T]{r: r, jd: json.NewDecoder(r)}
}

func (r *BatchResultReader[T]) Read(out *BatchResult[T]) error {
	return r.jd.Decode(out)
}
