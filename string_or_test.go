package zhipu

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringOr(t *testing.T) {
	data := struct {
		Item *StringOr[int] `json:"item,omitempty"`
	}{}
	data.Item = &StringOr[int]{}
	data.Item.SetString("test")

	b, err := json.Marshal(data)
	require.NoError(t, err)
	require.Equal(t, `{"item":"test"}`, string(b))

	data.Item.SetValue(1)
	b, err = json.Marshal(data)
	require.NoError(t, err)
	require.Equal(t, `{"item":1}`, string(b))

	err = json.Unmarshal([]byte(`{"item":"test2"}`), &data)
	require.NoError(t, err)
	require.NotNil(t, data.Item.String)
	require.Nil(t, data.Item.Value)
	require.Equal(t, "test2", *data.Item.String)

	err = json.Unmarshal([]byte(`{"item":2}`), &data)
	require.NoError(t, err)
	require.Nil(t, data.Item.String)
	require.NotNil(t, data.Item.Value)
	require.Equal(t, 2, *data.Item.Value)
}
