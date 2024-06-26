package zhipu

// Ptr returns a pointer to the value passed in.
// Example:
//
//	web_search_enable = zhipu.Ptr(false)
func Ptr[T any](v T) *T {
	return &v
}

// M is a shorthand for map[string]any.
type M = map[string]any
