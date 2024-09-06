package utilities

// returns a pointer to the value provided
func ToPointer[T any](value T) *T {
	return &value
}
