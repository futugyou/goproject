package agentic

import (
	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/types"
)

type AgenticInput struct {
	RequestID string `json:"requestId"`
	AgentID   string `json:"agentId"`
	types.RunAgentInput
}
