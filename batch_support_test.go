package zhipu

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBatchFileWriter(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	buf := &bytes.Buffer{}

	w := NewBatchFileWriter(buf)
	err = w.Write("batch-1", client.ChatCompletionService("a").AddMessage(ChatCompletionMessage{
		Role: "user", Content: "hello",
	}))
	require.NoError(t, err)
	err = w.Write("batch-2", client.EmbeddingService("c").SetInput("whoa"))
	require.NoError(t, err)
	err = w.Write("batch-3", client.ImageGenerationService("d").SetPrompt("whoa"))
	require.NoError(t, err)

	require.Equal(t, `{"body":{"messages":[{"role":"user","content":"hello"}],"model":"a"},"custom_id":"batch-1","method":"POST","url":"/v4/chat/completions"}
{"body":{"input":"whoa","model":"c"},"custom_id":"batch-2","method":"POST","url":"/v4/embeddings"}
{"body":{"model":"d","prompt":"whoa"},"custom_id":"batch-3","method":"POST","url":"/v4/images/generations"}
`, buf.String())
}
