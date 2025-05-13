package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/text"
)

type MemoryAnswer struct {
	StreamState     *StreamState `json:"streamState,omitempty"`
	Question        *string      `json:"question,omitempty"`
	NoResult        bool         `json:"noResult"`
	NoResultReason  *string      `json:"noResultReason,omitempty"`
	Result          string       `json:"text"`
	TokenUsage      []TokenUsage `json:"tokenUsage,omitempty"`
	RelevantSources []Citation   `json:"relevantSources,omitempty"`
}

func (a *MemoryAnswer) ToJson(optimizeForStream bool) string {
	if a == nil {
		return "{}"
	}

	if !optimizeForStream || (a.StreamState != nil && *a.StreamState != StreamStateAppend) {
		if j, err := json.Marshal(a); err != nil {
			return "{}"
		} else {
			return string(j)
		}
	}

	if j, err := json.Marshal(a); err != nil {
		return "{}"
	} else {
		var cl MemoryAnswer
		if err = json.Unmarshal(j, &cl); err != nil {
			return "{}"
		}

		if cl.Question != nil && len(*cl.Question) == 0 {
			cl.Question = nil
		}

		if len(cl.RelevantSources) == 0 {
			cl.RelevantSources = nil
		}
		if j, err := json.Marshal(cl); err != nil {
			return "{}"
		} else {
			return string(j)
		}
	}
}

func (a *MemoryAnswer) ToString() string {
	if a == nil {
		return ""
	}

	if a.NoResult || len(a.RelevantSources) == 0 {
		return a.Result
	}

	var builder strings.Builder
	text.AppendStringLine(builder, a.Result)

	sources := map[string]string{}
	for _, x := range a.RelevantSources {
		if len(x.Partitions) > 0 {
			date := x.Partitions[0].LastUpdate.Format("20060102150405.9999999")
			key := x.Index + x.Link
			sources[key] = fmt.Sprintf("  - %s [%s]", x.SourceName, date)
		}
	}

	text.AppendStringLine(builder, "- Sources:")
	for _, v := range sources {
		text.AppendStringLine(builder, v)
	}

	return builder.String()
}

type StreamState string

var (
	StreamStateError  StreamState = "error"
	StreamStateReset  StreamState = "reset"
	StreamStateAppend StreamState = "append"
	StreamStateLast   StreamState = "last"
)
