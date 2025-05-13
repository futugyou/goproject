package dataformats

import (
	"strconv"
	"strings"
)

const metaSentencesAreComplete string = "completeSentences"
const metaPageNumber string = "pageNumber"

type Chunk struct {
	Number      int64             `json:"number"`
	Content     string            `json:"content"`
	Metadata    map[string]string `json:"metadata"`
	IsSeparator bool              `json:"-"`
}

func (c *Chunk) SentencesAreComplete() bool {
	if c == nil {
		return false
	}

	if v, ok := c.Metadata[metaSentencesAreComplete]; ok {
		if _, err := strconv.ParseBool(v); err == nil {
			return true
		}
	}

	return false
}

func (c *Chunk) PageNumber() int64 {
	if c == nil {
		return -1
	}

	if v, ok := c.Metadata[metaPageNumber]; ok {
		if vv, err := strconv.ParseInt(v, 10, 64); err == nil {
			return vv
		}
	}

	return -1
}

func NewChunk(number int64, content *string, metadata map[string]string) *Chunk {
	con := ""
	if content != nil {
		con = *content
	}
	if metadata == nil {
		metadata = make(map[string]string)
	}
	return &Chunk{
		Number:      number,
		Content:     con,
		Metadata:    metadata,
		IsSeparator: false,
	}
}

func NewChunkWithByte(number int64, content byte) *Chunk {
	con := string(content)
	return &Chunk{
		Number:      number,
		Content:     con,
		Metadata:    make(map[string]string),
		IsSeparator: false,
	}
}

func NewChunkWithStringBuilder(number int64, builder strings.Builder) *Chunk {
	con := builder.String()
	return &Chunk{
		Number:      number,
		Content:     con,
		Metadata:    make(map[string]string),
		IsSeparator: false,
	}
}

func ChunkMeta(sentencesAreComplete *bool, pageNumber *int64) map[string]string {
	metadata := make(map[string]string)
	if sentencesAreComplete != nil {
		metadata[metaSentencesAreComplete] = strconv.FormatBool(*sentencesAreComplete)
	}
	if pageNumber != nil {
		metadata[metaPageNumber] = strconv.FormatInt(*pageNumber, 10)
	}
	return metadata
}
