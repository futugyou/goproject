package tiktoken

import (
	"context"
	"fmt"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/pkoukk/tiktoken-go"
)

type CL100KTokenizer struct {
	encoder *tiktoken.Tiktoken
}

func NewCL100KTokenizer() (*CL100KTokenizer, error) {
	enc, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return nil, err
	}
	return &CL100KTokenizer{encoder: enc}, nil
}

func (t *CL100KTokenizer) CountTokens(ctx context.Context, text string) int64 {
	if t == nil {
		return 0
	}
	return int64(len(t.encoder.Encode(text, nil, nil)))
}

func (t *CL100KTokenizer) GetTokens(ctx context.Context, text string) []string {
	if t == nil {
		return []string{}
	}
	tokenIDs := t.encoder.Encode(text, nil, nil)
	tokens := make([]string, len(tokenIDs))
	for i, id := range tokenIDs {
		tokens[i] = fmt.Sprintf("%d", id)
	}
	return tokens
}

var _ ai.ITextTokenizer = (*CL100KTokenizer)(nil)
