package zhipu

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileCreateServiceFineTune(t *testing.T) {
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

func TestFileCreateServiceKnowledge(t *testing.T) {
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
}
