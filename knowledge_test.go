package zhipu

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKnowledgeCapacity(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.KnowledgeCapacity()
	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.Total.Length)
	require.NotEmpty(t, res.Total.WordNum)
}

func TestKnowledgeServiceAll(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)

	s := client.KnowledgeCreate()
	s.SetName("test")
	s.SetDescription("test description")
	s.SetEmbeddingID(KnowledgeEmbeddingIDEmbedding2)

	res, err := s.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.ID)

	s2 := client.KnowledgeList()
	res2, err := s2.Do(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res2.List)
	require.Equal(t, res.ID, res2.List[0].ID)

	s3 := client.KnowledgeEdit(res.ID)
	s3.SetDescription("test description 2")
	s3.SetName("test 2")
	s3.SetEmbeddingID(KnowledgeEmbeddingIDEmbedding2)
	err = s3.Do(context.Background())
	require.NoError(t, err)

	s4 := client.KnowledgeDelete(res.ID)
	err = s4.Do(context.Background())
	require.NoError(t, err)
}
