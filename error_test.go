package zhipu

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIError(t *testing.T) {
	err := APIError{
		Code:    "code",
		Message: "message",
	}
	require.Equal(t, "message", err.Error())
	require.Equal(t, "code", GetAPIErrorCode(err))
	require.Equal(t, "message", GetAPIErrorMessage(err))
}

func TestAPIErrorResponse(t *testing.T) {
	err := APIErrorResponse{
		APIError: APIError{
			Code:    "code",
			Message: "message",
		},
	}
	require.Equal(t, "message", err.Error())
	require.Equal(t, "code", GetAPIErrorCode(err))
	require.Equal(t, "message", GetAPIErrorMessage(err))
}

func TestAPIErrorResponseFromDoc(t *testing.T) {
	var res APIErrorResponse
	err := json.Unmarshal([]byte(`{"error":{"code":"1002","message":"Authorization Token非法，请确认Authorization Token正确传递。"}}`), &res)
	require.NoError(t, err)
	require.Equal(t, "1002", res.Code)
	require.Equal(t, "1002", GetAPIErrorCode(res))
}
