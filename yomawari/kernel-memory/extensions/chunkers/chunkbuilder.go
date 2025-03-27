package chunkers

import "strings"

type ChunkBuilder struct {
	FullContent  strings.Builder
	NextSentence strings.Builder
}
