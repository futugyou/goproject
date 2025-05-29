package protocol

type ContextInclusion string

const (
	ContextInclusionNone       ContextInclusion = "none"
	ContextInclusionThisServer ContextInclusion = "thisServer"
	ContextInclusionAllServers ContextInclusion = "allServers"
)
