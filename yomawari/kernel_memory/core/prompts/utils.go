package prompts

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
)

var (
	tagsRegex     = regexp.MustCompile(`\{\{\$tags\[(.*?)\]\}\}`)
	metadataRegex = regexp.MustCompile(`\{\{\$meta\[(.*?)\]\}\}`)
)

func RenderFactTemplate(
	template, factContent string,
	source, relevance, recordId string,
	tags *models.TagCollection,
	metadata map[string]interface{},
) string {
	result := strings.ReplaceAll(template, "{{$source}}", source)
	result = strings.ReplaceAll(result, "{{$relevance}}", relevance)
	result = strings.ReplaceAll(result, "{{$memoryId}}", recordId)

	// {{$tag[X]}}
	for {
		matches := tagsRegex.FindStringSubmatch(result)
		if matches == nil {
			break
		}
		tagName := matches[1]
		if tags == nil {
			return "-"
		}

		tagValues, exists := tags.Get(tagName)
		replacement := "-"
		if exists {
			switch len(tagValues) {
			case 1:
				replacement = tagValues[0]
			case 0:
				replacement = "-"
			default:
				replacement = "[" + strings.Join(tagValues, ", ") + "]"
			}
		}
		result = tagsRegex.ReplaceAllString(result, replacement)
	}

	// {{$tags}}
	var tagList []string
	for key, values := range tags.GetData() {
		tagList = append(tagList, fmt.Sprintf("%s=[%s]", key, strings.Join(values, ", ")))
	}
	result = strings.ReplaceAll(result, "{{$tags}}", strings.Join(tagList, ", "))

	// {{$meta[X]}}
	for {
		matches := metadataRegex.FindStringSubmatch(result)
		if matches == nil {
			break
		}
		metaKey := matches[1]
		metaValue, exists := metadata[metaKey]
		replacement := "-"
		if exists {
			replacement = fmt.Sprintf("%v", metaValue)
		}
		result = metadataRegex.ReplaceAllString(result, replacement)
	}

	// {{$content}}
	result = strings.ReplaceAll(result, "{{$content}}", factContent)

	return result
}
