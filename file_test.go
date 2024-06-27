package zhipu

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileServiceFineTune(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.FileCreateService(FilePurposeFineTune)
	s.SetLocalFile(filepath.Join("testdata", "test-file.jsonl"))

	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotZero(t, res.Bytes)
	require.NotZero(t, res.CreatedAt)
	require.NotEmpty(t, res.ID)
}

func TestFileServiceKnowledge(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.FileCreateService(FilePurposeRetrieval)
	s.SetKnowledgeID(os.Getenv("TEST_KNOWLEDGE_ID"))
	s.SetLocalFile(filepath.Join("testdata", "test-file.txt"))

	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.SuccessInfos)
	require.NotEmpty(t, res.SuccessInfos[0].DocumentID)
	require.NotEmpty(t, res.SuccessInfos[0].Filename)

	documentID := res.SuccessInfos[0].DocumentID

	res2, err := client.FileGetService(documentID).Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res2.ID)

	err = client.FileEditService(documentID).SetKnowledgeType(KnowledgeTypeCustom).Do(context.Background())
	require.True(t, err == nil || GetAPIErrorCode(err) == "10019")

	err = client.FileDeleteService(res.SuccessInfos[0].DocumentID).Do(context.Background())
	require.True(t, err == nil || GetAPIErrorCode(err) == "10019")
}

func TestFileListServiceKnowledge(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.FileListService(FilePurposeRetrieval).SetKnowledgeID(os.Getenv("TEST_KNOWLEDGE_ID"))
	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.List)
}

func TestFileListServiceFineTune(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.FileListService(FilePurposeFineTune)
	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.Data)
}
