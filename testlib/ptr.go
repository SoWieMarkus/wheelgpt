package testlib

// Init something as a pointer.
func Ptr[T any](v T) *T { return &v }
