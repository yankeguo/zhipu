package zhipu

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageGenerationService(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.ImageGeneration("cogview-3")
	s.SetPrompt("一只可爱的小猫")

	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.Data)
	t.Log(res.Data[0].URL)
}
