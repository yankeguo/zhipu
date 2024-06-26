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

func TestChatCompletionServiceStream(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	var content string

	s := client.ChatCompletionService("glm-4-flash").AddMessage(ChatCompletionMessage{
		Role:    "user",
		Content: "你好呀",
	}).SetStreamHandler(func(chunk ChatCompletionResponse) error {
		content += chunk.Choices[0].Delta.Content
		return nil
	})
	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.Choices)
	choice := res.Choices[0]
	require.Equal(t, FinishReasonStop, choice.FinishReason)
	require.NotEmpty(t, choice.Message.Content)
	require.Equal(t, content, choice.Message.Content)
}
