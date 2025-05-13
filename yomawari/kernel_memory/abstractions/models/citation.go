package models

import "time"

type Citation struct {
	Link              string      `json:"link"`
	Index             string      `json:"index"`
	DocumentId        string      `json:"documentId"`
	FileId            string      `json:"fileId"`
	SourceContentType string      `json:"sourceContentType"`
	SourceName        string      `json:"sourceName"`
	SourceUrl         *string     `json:"sourceUrl,omitempty"`
	Partitions        []Partition `json:"partitions"`
}

type Partition struct {
	Text            string         `json:"text"`
	Relevance       float64        `json:"relevance"`
	PartitionNumber int64          `json:"partitionNumber"`
	SectionNumber   int64          `json:"sectionNumber"`
	LastUpdate      time.Time      `json:"lastUpdate"`
	Tags            *TagCollection `json:"tags"`
}
