package models

import "time"

type DataPipelineStatus struct {
	Completed      bool           `json:"completed"`
	Empty          bool           `json:"empty"`
	Index          string         `json:"index"`
	DocumentId     string         `json:"document_id"`
	Tags           *TagCollection `json:"tags"`
	Creation       time.Time      `json:"creation"`
	LastUpdate     time.Time      `json:"last_update"`
	Steps          []string       `json:"steps"`
	RemainingSteps []string       `json:"remaining_steps"`
	CompletedSteps []string       `json:"completed_steps"`
}
