package configuration

import "fmt"

type TextPartitioningOptions struct {
	MaxTokensPerParagraph int64
	OverlappingTokens     int64
}

func NewTextPartitioningOptions() *TextPartitioningOptions {
	return &TextPartitioningOptions{
		MaxTokensPerParagraph: 1000,
		OverlappingTokens:     100,
	}
}

func (o *TextPartitioningOptions) Validate() error {
	if o == nil {
		return fmt.Errorf("nil options")
	}
	if o.MaxTokensPerParagraph < 1 {
		return fmt.Errorf("text partitioning: MaxTokensPerParagraph cannot be less than 1")
	}
	if o.OverlappingTokens < 0 {
		return fmt.Errorf("text partitioning: OverlappingTokens cannot be less than 1")
	}
	if o.OverlappingTokens >= o.MaxTokensPerParagraph {
		return fmt.Errorf("text partitioning: OverlappingTokens cannot be greater than MaxTokensPerParagraph")
	}
	return nil
}
