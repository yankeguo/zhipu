package zhipu

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChatCompletionService(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.ChatCompletion("glm-4-flash")
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

	s := client.ChatCompletion("charglm-3")
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
	require.Contains(t, []string{FinishReasonStop, FinishReasonStopSequence}, choice.FinishReason)
	require.NotEmpty(t, choice.Message.Content)
}

func TestChatCompletionServiceAllToolsCodeInterpreter(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.ChatCompletion("GLM-4-AllTools")
	s.AddMessage(ChatCompletionMultiMessage{
		Role: "user",
		Content: []ChatCompletionMultiContent{
			{
				Type: "text",
				Text: "计算[5,10,20,700,99,310,978,100]的平均值和方差。",
			},
		},
	})
	s.AddTool(ChatCompletionToolCodeInterpreter{
		Sandbox: Ptr(CodeInterpreterSandboxAuto),
	})

	foundInterpreterInput := false
	foundInterpreterOutput := false

	s.SetStreamHandler(func(chunk ChatCompletionResponse) error {
		for _, c := range chunk.Choices {
			for _, tc := range c.Delta.ToolCalls {
				if tc.Type == ToolTypeCodeInterpreter && tc.CodeInterpreter != nil {
					if tc.CodeInterpreter.Input != "" {
						foundInterpreterInput = true
					}
					if len(tc.CodeInterpreter.Outputs) > 0 {
						foundInterpreterOutput = true
					}
				}
			}
		}
		buf, _ := json.MarshalIndent(chunk, "", "  ")
		t.Log(string(buf))
		return nil
	})

	res, err := s.Do(context.Background())
	require.True(t, foundInterpreterInput)
	require.True(t, foundInterpreterOutput)
	require.NotNil(t, res)
	require.NoError(t, err)
}

func TestChatCompletionServiceAllToolsDrawingTool(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.ChatCompletion("GLM-4-AllTools")
	s.AddMessage(ChatCompletionMultiMessage{
		Role: "user",
		Content: []ChatCompletionMultiContent{
			{
				Type: "text",
				Text: "画一个正弦函数图像",
			},
		},
	})
	s.AddTool(ChatCompletionToolDrawingTool{})

	foundInput := false
	foundOutput := false
	outputImage := ""

	s.SetStreamHandler(func(chunk ChatCompletionResponse) error {
		for _, c := range chunk.Choices {
			for _, tc := range c.Delta.ToolCalls {
				if tc.Type == ToolTypeDrawingTool && tc.DrawingTool != nil {
					if tc.DrawingTool.Input != "" {
						foundInput = true
					}
					if len(tc.DrawingTool.Outputs) > 0 {
						foundOutput = true
					}
					for _, output := range tc.DrawingTool.Outputs {
						if output.Image != "" {
							outputImage = output.Image
						}
					}
				}
			}
		}
		buf, _ := json.MarshalIndent(chunk, "", "  ")
		t.Log(string(buf))
		return nil
	})

	res, err := s.Do(context.Background())
	require.True(t, foundInput)
	require.True(t, foundOutput)
	require.NotEmpty(t, outputImage)
	t.Log(outputImage)
	require.NotNil(t, res)
	require.NoError(t, err)
}

func TestChatCompletionServiceAllToolsWebBrowser(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.ChatCompletion("GLM-4-AllTools")
	s.AddMessage(ChatCompletionMultiMessage{
		Role: "user",
		Content: []ChatCompletionMultiContent{
			{
				Type: "text",
				Text: "搜索下本周深圳天气如何",
			},
		},
	})
	s.AddTool(ChatCompletionToolWebBrowser{})

	foundInput := false
	foundOutput := false
	outputContent := ""

	s.SetStreamHandler(func(chunk ChatCompletionResponse) error {
		for _, c := range chunk.Choices {
			for _, tc := range c.Delta.ToolCalls {
				if tc.Type == ToolTypeWebBrowser && tc.WebBrowser != nil {
					if tc.WebBrowser.Input != "" {
						foundInput = true
					}
					if len(tc.WebBrowser.Outputs) > 0 {
						foundOutput = true
					}
					for _, output := range tc.WebBrowser.Outputs {
						if output.Content != "" {
							outputContent = output.Content
						}
					}
				}
			}
		}
		buf, _ := json.MarshalIndent(chunk, "", "  ")
		t.Log(string(buf))
		return nil
	})

	res, err := s.Do(context.Background())
	require.True(t, foundInput)
	require.True(t, foundOutput)
	require.NotEmpty(t, outputContent)
	t.Log(outputContent)
	require.NotNil(t, res)
	require.NoError(t, err)
}

func TestChatCompletionServiceStream(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	var content string

	s := client.ChatCompletion("glm-4-flash").AddMessage(ChatCompletionMessage{
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

	s := client.ChatCompletion("glm-4v")
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
