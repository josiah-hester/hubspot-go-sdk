// Package ptr provides helper functions for creating pointers to values.
package ptr

// To returns a pointer to the given value. Useful for constructing
// request structs with optional fields.
func To[T any](v T) *T {
	return &v
}
