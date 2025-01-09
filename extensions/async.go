package extensions

import "context"

func HandleAsync[T any](ctx context.Context, resCh <-chan T, errCh <-chan error) (T, error) {
	var zero T
	select {
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return zero, err
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}

func HandleErrorAsync(ctx context.Context, errCh <-chan error) error {
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
