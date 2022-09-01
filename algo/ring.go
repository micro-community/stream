package algo

// Ring for anytype
type Ring[T any] struct {
	next, prev *Ring[T]
	Value      T // for use by client; untouched by this library
}
