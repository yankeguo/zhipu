package zhipu

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientR(t *testing.T) {
	c, err := NewClient()
	require.NoError(t, err)
	// the only free api is to list fine-tuning jobs
	res, err := c.request(context.Background()).Get("fine_tuning/jobs")
	require.NoError(t, err)
	require.True(t, res.IsSuccess())
}
