package zhipu

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

type APIErrorResponse struct {
	APIError `json:"error"`
}

func (e APIErrorResponse) Error() string {
	return e.APIError.Error()
}

// GetAPIErrorCode returns the error code of an API error.
func GetAPIErrorCode(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(APIError); ok {
		return e.Code
	}
	if e, ok := err.(APIErrorResponse); ok {
		return e.Code
	}
	if e, ok := err.(*APIError); ok && e != nil {
		return e.Code
	}
	if e, ok := err.(*APIErrorResponse); ok && e != nil {
		return e.Code
	}
	return ""
}

// GetAPIErrorMessage returns the error message of an API error.
func GetAPIErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(APIError); ok {
		return e.Message
	}
	if e, ok := err.(APIErrorResponse); ok {
		return e.Message
	}
	if e, ok := err.(*APIError); ok && e != nil {
		return e.Message
	}
	if e, ok := err.(*APIErrorResponse); ok && e != nil {
		return e.Message
	}
	return err.Error()
}
