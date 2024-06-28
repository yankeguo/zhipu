package zhipu

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBatchServiceAll(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	buf := &bytes.Buffer{}

	bfw := NewBatchFileWriter(buf)
	err = bfw.Write("batch_1", client.ChatCompletion("glm-4-flash").AddMessage(ChatCompletionMessage{
		Role: RoleUser, Content: "你好呀",
	}))
	require.NoError(t, err)
	err = bfw.Write("batch_2", client.ChatCompletion("glm-4-flash").AddMessage(ChatCompletionMessage{
		Role: RoleUser, Content: "你叫什么名字",
	}))
	require.NoError(t, err)

	res, err := client.FileCreate(FilePurposeBatch).SetFile(bytes.NewReader(buf.Bytes()), "batch.jsonl").Do(context.Background())
	require.NoError(t, err)

	fileID := res.FileCreateFineTuneResponse.ID
	require.NotEmpty(t, fileID)

	res1, err := client.BatchCreate().
		SetInputFileID(fileID).
		SetCompletionWindow(BatchCompletionWindow24h).
		SetEndpoint(BatchEndpointV4ChatCompletions).Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res1.ID)

	res2, err := client.BatchGet(res1.ID).Do(context.Background())
	require.NoError(t, err)
	require.Equal(t, res2.ID, res1.ID)

	res3, err := client.BatchList().Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res3.Data)

	err = client.BatchCancel(res1.ID).Do(context.Background())
	require.NoError(t, err)
}

func TestBatchListService(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	res, err := client.BatchList().Do(context.Background())
	require.NoError(t, err)
	t.Log(res)
}
