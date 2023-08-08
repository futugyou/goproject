package main

type CloneInfo struct {
	SourceOwner  string
	SourceBranch string
	SourceName   string
	SourceToken  string

	DestOwner  string
	DestBranch string
	DestName   string
	DestToken  string
}
