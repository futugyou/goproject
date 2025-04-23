package tasks

type Resolver[T any] interface {
	SetResult(value T)
	SetError(err error)
	SetCanceled()
	TrySetResult(value T) bool
	TrySetError(err error) bool
	TrySetCanceled() bool
}
