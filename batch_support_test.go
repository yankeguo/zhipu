package zhipu

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBatchFileWriter(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	buf := &bytes.Buffer{}

	w := NewBatchFileWriter(buf)
	err = w.Write("batch-1", client.ChatCompletion("a").AddMessage(ChatCompletionMessage{
		Role: "user", Content: "hello",
	}))
	require.NoError(t, err)
	err = w.Write("batch-2", client.Embedding("c").SetInput("whoa"))
	require.NoError(t, err)
	err = w.Write("batch-3", client.ImageGeneration("d").SetPrompt("whoa"))
	require.NoError(t, err)

	require.Equal(t, `{"body":{"messages":[{"role":"user","content":"hello"}],"model":"a"},"custom_id":"batch-1","method":"POST","url":"/v4/chat/completions"}
{"body":{"input":"whoa","model":"c"},"custom_id":"batch-2","method":"POST","url":"/v4/embeddings"}
{"body":{"model":"d","prompt":"whoa"},"custom_id":"batch-3","method":"POST","url":"/v4/images/generations"}
`, buf.String())
}

func TestBatchResultReader(t *testing.T) {
	result := `
	{"response":{"status_code":200,"body":{"created":1715959701,"usage":{"completion_tokens":26,"prompt_tokens":89,"total_tokens":115},"model":"glm-4","id":"8668357533850320547","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n  \"分类标签\": \"负面\",\n  \"特定问题标注\": \"订单处理慢\"\n}\n'''"}}],"request_id":"615-request-1"}},"custom_id":"request-1","id":"batch_1791490810192076800"}
{"response":{"status_code":200,"body":{"created":1715959701,"usage":{"completion_tokens":22,"prompt_tokens":94,"total_tokens":116},"model":"glm-4","id":"8668368425887509080","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n \"分类标签\": \"负面\",\n \"特定问题标注\": \"产品缺陷\"\n}\n'''"}}],"request_id":"616-request-2"}},"custom_id":"request-2","id":"batch_1791490810192076800"}
{"response":{"status_code":200,"body":{"created":1715959701,"usage":{"completion_tokens":25,"prompt_tokens":86,"total_tokens":111},"model":"glm-4","id":"8668355815863214980","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n  \"分类标签\": \"正面\",\n  \"特定问题标注\": \"性价比\"\n}\n'''"}}],"request_id":"617-request-3"}},"custom_id":"request-3","id":"batch_1791490810192076800"}
{"response":{"status_code":200,"body":{"created":1715959701,"usage":{"completion_tokens":28,"prompt_tokens":89,"total_tokens":117},"model":"glm-4","id":"8668355815863214981","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n  \"分类标签\": \"负面\",\n  \"特定问题标注\": \"说明文档不清晰\"\n}\n'''"}}],"request_id":"618-request-4"}},"custom_id":"request-4","id":"batch_1791490810192076800"}

{"response":{"status_code":200,"body":{"created":1715959701,"usage":{"completion_tokens":26,"prompt_tokens":88,"total_tokens":114},"model":"glm-4","id":"8668357533850320546","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n  \"分类标签\": \"中性\",\n  \"特定问题标注\": \"价格问题\"\n}\n'''"}}],"request_id":"619-request-5"}},"custom_id":"request-5","id":"batch_1791490810192076800"}

{"response":{"status_code":200,"body":{"created":1715959701,"usage":{"completion_tokens":26,"prompt_tokens":90,"total_tokens":116},"model":"glm-4","id":"8668356159460662846","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n  \"分类标签\": \"负面\",\n  \"特定问题标注\": \"配送延迟\"\n}\n'''"}}],"request_id":"620-request-6"}},"custom_id":"request-6","id":"batch_1791490810192076800"}


{"response":{"status_code":200,"body":{"created":1715959701,"usage":{"completion_tokens":27,"prompt_tokens":88,"total_tokens":115},"model":"glm-4","id":"8668357671289274638","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n  \"分类标签\": \"负面\",\n  \"特定问题标注\": \"产品描述不符\"\n}\n'''"}}],"request_id":"621-request-7"}},"custom_id":"request-7","id":"batch_1791490810192076800"}
{"response":{"status_code":200,"body":{"created":1715959702,"usage":{"completion_tokens":26,"prompt_tokens":87,"total_tokens":113},"model":"glm-4","id":"8668355644064514872","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n  \"分类标签\": \"正面\",\n  \"特定问题标注\": \"客服态度\"\n}\n'''"}}],"request_id":"622-request-8"}},"custom_id":"request-8","id":"batch_1791490810192076800"}
    {"response":{"status_code":200,"body":{"created":1715959701,"usage":{"completion_tokens":29,"prompt_tokens":90,"total_tokens":119},"model":"glm-4","id":"8668357671289274639","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n  \"分类标签\": \"负面\",\n  \"特定问题标注\": \"包装问题, 产品损坏\"\n}\n'''"}}],"request_id":"623-request-9"}},"custom_id":"request-9","id":"batch_1791490810192076800"}
{"response":{"status_code":200,"body":{"created":1715959701,"usage":{"completion_tokens":27,"prompt_tokens":87,"total_tokens":114},"model":"glm-4","id":"8668355644064514871","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"'''json\n{\n  \"分类标签\": \"正面\",\n  \"特定问题标注\": \"产品描述不符\"\n}\n'''"}}],"request_id":"624-request-10"}},"custom_id":"request-10","id":"batch_1791490810192076800"}
`

	brr := NewBatchResultReader[ChatCompletionResponse](bytes.NewReader([]byte(result)))

	var count int

	for {
		var res BatchResult[ChatCompletionResponse]

		err := brr.Read(&res)

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			require.Equal(t, 10, count)
			require.NoError(t, err)
			break
		}

		require.Equal(t, 200, res.Response.StatusCode)

		count++
	}
}
