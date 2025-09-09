package models

type FlowGraph struct {
	Nodes    []Node   `json:"nodes"`
	Edges    []Edge   `json:"edges"`
	Viewport Viewport `json:"viewport"`
}

type Edge struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Source       string         `json:"source"`
	SourceHandle string         `json:"sourceHandle"`
	Target       string         `json:"target"`
	TargetHandle string         `json:"targetHandle"`
	Data         map[string]any `json:"data,omitempty"`
	Style        map[string]any `json:"style,omitempty"`
	MarkerEnd    map[string]any `json:"markerEnd,omitempty"`
}

type Node struct {
	ID               string         `json:"id"`
	Type             string         `json:"type"`
	Position         PositionClass  `json:"position"`
	Data             map[string]any `json:"data"`
	PositionAbsolute PositionClass  `json:"positionAbsolute"`
	Selected         bool           `json:"selected"`
	Dragging         bool           `json:"dragging"`
	Style            map[string]any `json:"style,omitempty"`
	ClassName        string         `json:"className"`
	Resizing         *bool          `json:"resizing,omitempty"`
	ParentNode       *string        `json:"parentNode,omitempty"`
	Width            int64          `json:"width"`
	Height           int64          `json:"height"`
}

type PositionClass struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Viewport struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Zoom float64 `json:"zoom"`
}
