package zhipu

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmbeddingService(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	service := client.Embedding("embedding-2")

	resp, err := service.SetInput("你好").Do(context.Background())
	require.NoError(t, err)
	require.NotZero(t, resp.Usage.TotalTokens)
	require.NotEmpty(t, resp.Data)
	require.NotEmpty(t, resp.Data[0].Embedding)
}
