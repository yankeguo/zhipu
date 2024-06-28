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
		Role:    RoleUser,
		Content: "你好呀",
	})
	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.Choices)
	choice := res.Choices[0]
	require.Equal(t, FinishReasonStop, choice.FinishReason)
	require.NotEmpty(t, choice.Message.Content)
}

func TestChatCompletionServiceCharGLM(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.ChatCompletionService("charglm-3")
	s.SetMeta(
		ChatCompletionMeta{
			UserName: "啵酱",
			UserInfo: "啵酱是小少爷",
			BotName:  "塞巴斯酱",
			BotInfo:  "塞巴斯酱是一个冷酷的恶魔管家",
		},
	).AddMessage(ChatCompletionMessage{
		Role:    RoleUser,
		Content: "早上好",
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
		Role:    RoleUser,
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

func TestChatCompletionServiceVision(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.ChatCompletionService("glm-4v")
	s.AddMessage(ChatCompletionMultiMessage{
		Role: RoleUser,
		Content: []ChatCompletionMultiContent{
			{
				Type: MultiContentTypeText,
				Text: "图里有什么",
			},
			{
				Type: MultiContentTypeImageURL,
				ImageURL: &URLItem{
					URL: "https://img1.baidu.com/it/u=1369931113,3388870256&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1703696400&t=f3028c7a1dca43a080aeb8239f09cc2f",
				},
			},
		},
	})
	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.Choices)
	require.NotZero(t, res.Usage.CompletionTokens)
	choice := res.Choices[0]
	require.Equal(t, FinishReasonStop, choice.FinishReason)
	require.NotEmpty(t, choice.Message.Content)
}
