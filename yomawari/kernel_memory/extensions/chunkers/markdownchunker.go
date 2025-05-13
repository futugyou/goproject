package chunkers

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/dataformats"

	aitext "github.com/futugyou/yomawari/kernel_memory/abstractions/text"
)

type MarkDownChunkerOptions struct {
	MaxTokensPerChunk int
	Overlap           int
	ChunkHeader       *string
}

var explicitSeparatorsMD *SeparatorTrie = NewSeparatorTrie([]string{
	".\n\n",
	"!\n\n",
	"!!\n\n",
	"!!!\n\n",
	"?\n\n",
	"??\n\n",
	"???\n\n",
	"\n\n",
	"\n#",
	"\n##",
	"\n###",
	"\n####",
	"\n#####",
	"\n---",
})

// Prioritized list of characters to split inside a sentence.
var potentialSeparatorsMD *SeparatorTrie = NewSeparatorTrie([]string{
	"\n> ",
	"\n>- ",
	"\n>* ",
	"\n1. ",
	"\n2. ",
	"\n3. ",
	"\n4. ",
	"\n5. ",
	"\n6. ",
	"\n7. ",
	"\n8. ",
	"\n9. ",
	"\n10. ",
	"\n```",
})

// Prioritized list of characters to split inside a sentence when other splits are not found.
var weakSeparators1MD *SeparatorTrie = NewSeparatorTrie([]string{
	"![",
	"[",
	"| ",
	" |\n",
	"-|\n",
	"\n: ",
})

// Prioritized list of characters to split inside a sentence when other splits are not found.
var weakSeparators2MD *SeparatorTrie = NewSeparatorTrie([]string{
	// Symbol + space
	". ", ".\t", ".\n",
	"? ", "?\t", "?\n",
	"! ", "!\t", "!\n",
	"⁉ ", "⁉\t", "⁉\n",
	"⁈ ", "⁈\t", "⁈\n",
	"⁇ ", "⁇\t", "⁇\n",
	"… ", "…\t", "…\n",
	// Multi-char separators without space, ordered by length
	"!!!!", "????", "!!!", "???", "?!?", "!?!", "!?", "?!", "!!", "??", "....", "...", "..",
	// 1 char separators without space
	".", "?", "!", "⁉", "⁈", "⁇", "…",
})

// Prioritized list of characters to split inside a sentence when other splits are not found.
var weakSeparators3MD *SeparatorTrie = NewSeparatorTrie([]string{
	"; ", ";\t", ";\n", ";",
	"} ", "}\t", "}\n", "}",
	") ", ")\t", ")\n",
	"] ", "]\t", "]\n",
	")", "]",
	": ", ":",
	", ", ",",
	"\n",
})

type MarkDownChunker struct {
	tokenizer ai.ITextTokenizer
}

// NewMarkDownChunker creates a new MarkDownChunker instance
func NewMarkDownChunker(tokenizer ai.ITextTokenizer) *MarkDownChunker {
	if tokenizer == nil {
		// Default to CL100KTokenizer if not provided
		// You'll need to implement this or provide it
		panic("default tokenizer not implemented")
	}
	return &MarkDownChunker{tokenizer: tokenizer}
}

// Split splits text into chunks with a maximum token count
func (c *MarkDownChunker) Split(text string, maxTokensPerChunk int) []string {
	return c.SplitWithOptions(text, MarkDownChunkerOptions{MaxTokensPerChunk: maxTokensPerChunk})
}

// SplitWithOptions splits text into chunks with advanced options
func (c *MarkDownChunker) SplitWithOptions(text string, options MarkDownChunkerOptions) []string {
	if text == "" {
		return []string{}
	}

	// Clean up text (normalize newlines)
	text = aitext.NormalizeNewlines(text, true)

	// Calculate chunk sizes with headers and overlaps
	maxChunk1Size := options.MaxTokensPerChunk - c.TokenCount(options.ChunkHeader)
	maxChunkNSize := options.MaxTokensPerChunk - c.TokenCount(options.ChunkHeader) - options.Overlap
	maxChunk1Size = int(math.Max(MinChunkSize, float64(maxChunk1Size)))
	maxChunkNSize = int(math.Max(MinChunkSize, float64(maxChunkNSize)))

	// Recursive chunking
	firstChunkDone := false
	chunks := c.recursiveSplit(text, maxChunk1Size, maxChunkNSize, ExplicitSeparator, &firstChunkDone)

	// Add overlapping tokens if needed
	if options.Overlap > 0 && len(chunks) > 1 {
		newChunks := []string{chunks[0]}
		for i := 1; i < len(chunks); i++ {
			prevTokens := c.tokenizer.GetTokens(context.Background(), chunks[i-1])
			start := int(math.Max(0, float64(len(prevTokens)-options.Overlap)))
			overlapTokens := prevTokens[start:]
			newChunks = append(newChunks, strings.Join(overlapTokens, "")+chunks[i])
		}
		chunks = newChunks
	}

	// Add header to each chunk
	if options.ChunkHeader != nil && *options.ChunkHeader != "" {
		for i := range chunks {
			chunks[i] = *options.ChunkHeader + chunks[i]
		}
	}

	return chunks
}

// recursiveSplit splits text recursively using different separator types
func (c *MarkDownChunker) recursiveSplit(
	text string,
	maxChunk1Size int,
	maxChunkNSize int,
	separatorType SeparatorTypes,
	firstChunkDone *bool,
) []string {
	if text == "" {
		return []string{}
	}

	maxChunkSize := maxChunkNSize
	if !*firstChunkDone {
		maxChunkSize = maxChunk1Size
	}

	if c.TokenCount(&text) <= maxChunkSize {
		return []string{text}
	}

	var fragments []dataformats.Chunk
	switch separatorType {
	case ExplicitSeparator:
		fragments = SplitToFragments(text, explicitSeparatorsMD)
	case PotentialSeparator:
		fragments = SplitToFragments(text, potentialSeparatorsMD)
	case WeakSeparator1:
		fragments = SplitToFragments(text, weakSeparators1MD)
	case WeakSeparator2:
		fragments = SplitToFragments(text, weakSeparators2MD)
	case WeakSeparator3:
		fragments = SplitToFragments(text, weakSeparators3MD)
	case NotASeparator:
		fragments = SplitToFragments(text, nil)
	default:
		panic(fmt.Sprintf("unknown separator type: %v", separatorType))
	}

	return c.generateChunks(fragments, maxChunk1Size, maxChunkNSize, separatorType, firstChunkDone)
}

func (c *MarkDownChunker) TokenCount(input *string) int {
	if input == nil {
		return 0
	}
	return int(c.tokenizer.CountTokens(context.Background(), *input))
}

// generateChunks generates chunks from fragments
func (c *MarkDownChunker) generateChunks(
	fragments []dataformats.Chunk,
	maxChunk1Size int,
	maxChunkNSize int,
	separatorType SeparatorTypes,
	firstChunkDone *bool,
) []string {
	if len(fragments) == 0 {
		return []string{}
	}

	var chunks []string
	builder := ChunkBuilder{
		FullContent:  strings.Builder{},
		NextSentence: strings.Builder{},
	}
	var maxChunkSize int
	for _, fragment := range fragments {
		builder.NextSentence.WriteString(fragment.Content)

		if !fragment.IsSeparator {
			continue
		}

		nextSentence := builder.NextSentence.String()
		nextSentenceSize := c.TokenCount(&nextSentence)
		maxChunkSize = maxChunkNSize
		if !*firstChunkDone {
			maxChunkSize = maxChunk1Size
		}

		var state int
		if builder.FullContent.Len() == 0 {
			if nextSentenceSize <= maxChunkSize {
				state = 1
			} else {
				state = 2
			}
		} else {
			if nextSentenceSize <= maxChunkSize {
				state = 3
			} else {
				state = 4
			}
		}

		switch state {
		case 1:
			builder.FullContent.WriteString(nextSentence)
			builder.NextSentence.Reset()
		case 2:
			moreChunks := c.recursiveSplit(nextSentence, maxChunk1Size, maxChunkNSize, NextSeparatorType(separatorType), firstChunkDone)
			chunks = append(chunks, moreChunks[:len(moreChunks)-1]...)
			builder.NextSentence.Reset()
			builder.NextSentence.WriteString(moreChunks[len(moreChunks)-1])
		case 3:
			chunkPlusSentence := builder.FullContent.String() + builder.NextSentence.String()
			if c.TokenCount(&chunkPlusSentence) <= maxChunkSize {
				builder.FullContent.WriteString(builder.NextSentence.String())
			} else {
				AddChunkString(chunks, builder.FullContent.String(), firstChunkDone)
				builder.FullContent.Reset()
				builder.FullContent.WriteString(builder.NextSentence.String())
			}
			builder.NextSentence.Reset()
		case 4:
			AddChunkString(chunks, builder.FullContent.String(), firstChunkDone)
			moreChunks := c.recursiveSplit(nextSentence, maxChunk1Size, maxChunkNSize, NextSeparatorType(separatorType), firstChunkDone)
			chunks = append(chunks, moreChunks[:len(moreChunks)-1]...)
			builder.NextSentence.Reset()
			builder.NextSentence.WriteString(moreChunks[len(moreChunks)-1])
		default:
			panic(fmt.Sprintf("invalid state: %d", state))
		}
	}

	// Handle remaining content
	fullSentenceLeft := builder.FullContent.String()
	nextSentenceLeft := builder.NextSentence.String()
	maxChunkSize = maxChunkNSize
	if !*firstChunkDone {
		maxChunkSize = maxChunk1Size
	}

	if len(fullSentenceLeft) > 0 || len(nextSentenceLeft) > 0 {
		combined := fullSentenceLeft + nextSentenceLeft
		if c.TokenCount(&combined) <= maxChunkSize {
			AddChunkString(chunks, combined, firstChunkDone)
		} else {
			if len(fullSentenceLeft) > 0 {
				AddChunkString(chunks, fullSentenceLeft, firstChunkDone)
			}
			if len(nextSentenceLeft) > 0 {
				if c.TokenCount(&nextSentenceLeft) < maxChunkSize {
					AddChunkString(chunks, nextSentenceLeft, firstChunkDone)
				} else {
					moreChunks := c.recursiveSplit(nextSentenceLeft, maxChunk1Size, maxChunkNSize, NextSeparatorType(separatorType), firstChunkDone)
					chunks = append(chunks, moreChunks...)
				}
			}
		}
	}

	return chunks
}
