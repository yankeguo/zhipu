package zhipu

// URLItem is a struct that contains a URL.
type URLItem struct {
	URL string `json:"url,omitempty"`
}

// IDItem is a struct that contains an ID.
type IDItem struct {
	ID string `json:"id,omitempty"`
}

// Ptr returns a pointer to the value passed in.
// Example:
//
//	web_search_enable = zhipu.Ptr(false)
func Ptr[T any](v T) *T {
	return &v
}

// M is a shorthand for map[string]any.
type M = map[string]any
