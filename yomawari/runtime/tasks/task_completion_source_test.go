package tasks

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTaskCompletionSource_ResultSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	tcs := NewTaskCompletionSource[string](ctx, cancel)

	go func() {
		time.Sleep(100 * time.Millisecond)
		tcs.SetResult("test-value")
	}()

	result, err := tcs.Result()
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if result != "test-value" {
		t.Fatalf("Expected 'test-value', got '%s'", result)
	}
}

func TestTaskCompletionSource_ResultTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	tcs := NewTaskCompletionSource[string](ctx, cancel)

	_, err := tcs.Result()
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected DeadlineExceeded, got %v", err)
	}
}

func TestTaskCompletionSource_CancelManually(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	tcs := NewTaskCompletionSource[string](ctx, cancel)

	go func() {
		time.Sleep(50 * time.Millisecond)
		tcs.Cancel()
	}()

	_, err := tcs.Result()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled, got %v", err)
	}
}
