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
		model, _ := m.Get(ctx, "id1")
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

func TestJwksGetWorks(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	m := NewMockIJwksRepository(ctrl)
	model := &token.JwkModel{}
	m.EXPECT().Get(ctx, "id1").Return(model, nil)
	if _, err := call(ctx, m); err != nil {
		t.Error("call failed:", err)
	}
}

func TestJwkModel_GetType(t *testing.T) {
	type fields struct {
		ID      string
		Payload string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := token.JwkModel{
				ID:      tt.fields.ID,
				Payload: tt.fields.Payload,
			}
			if got := j.GetType(); got != tt.want {
				t.Errorf("JwkModel.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}
