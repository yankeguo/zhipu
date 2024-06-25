package zhipu

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChatCompletionService(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.ChatCompletionService("glm-4-flash")
	s.AddMessage(ChatCompletionMessage{
		Role:    "user",
		Content: "你好呀",
	})
	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.Choices)
	choice := res.Choices[0]
	require.Equal(t, FinishReasonStop, choice.FinishReason)
	require.NotEmpty(t, choice.Message.Content)
}
