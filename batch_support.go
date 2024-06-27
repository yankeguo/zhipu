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
