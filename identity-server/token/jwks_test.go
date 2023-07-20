package token_test

import (
	"context"
	"testing"

	token "github.com/futugyousuzu/identity-server/token"
	"go.uber.org/mock/gomock"
)

func call(ctx context.Context, m token.IJwksRepository) (*token.JwkModel, error) {
	result := make(chan *token.JwkModel)
	go func() {
		model, _ := m.Get(ctx, "")
		result <- model
		close(result)
	}()
	select {
	case r := <-result:
		return r, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func TestJwksGetFails(t *testing.T) {
	t.Skip("Test is expected to fail, remove skip to trying running yourself.")
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	m := NewMockIJwksRepository(ctrl)
	if _, err := call(ctx, m); err != nil {
		t.Error("call failed:", err)
	}
}
