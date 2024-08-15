package zhipu

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVideoGeneration(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.VideoGeneration("cogvideox")
	s.SetPrompt("一只可爱的小猫")

	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.TaskStatus)
	require.NotEmpty(t, res.ID)
	t.Log(res.ID)

	for {
		res, err := client.AsyncResult(res.ID).Do(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, res.TaskStatus)
		if res.TaskStatus == VideoGenerationTaskStatusSuccess {
			require.NotEmpty(t, res.VideoResult)
			t.Log(res.VideoResult[0].URL)
			t.Log(res.VideoResult[0].CoverImageURL)
		}
		if res.TaskStatus != VideoGenerationTaskStatusProcessing {
			break
		}
		time.Sleep(time.Second * 5)
	}
}
