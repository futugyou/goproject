package models

type FlowGraph struct {
	Nodes    []Node   `json:"nodes" bson:"nodes"`
	Edges    []Edge   `json:"edges" bson:"edges"`
	Viewport Viewport `json:"viewport" bson:"viewport"`
}

type Edge struct {
	ID           string         `json:"id" bson:"id"`
	Type         string         `json:"type" bson:"type"`
	Source       string         `json:"source" bson:"source"`
	SourceHandle string         `json:"sourceHandle" bson:"sourceHandle"`
	Target       string         `json:"target" bson:"target"`
	TargetHandle string         `json:"targetHandle" bson:"targetHandle"`
	Data         map[string]any `json:"data,omitempty" bson:"data,omitempty"`
	Style        map[string]any `json:"style,omitempty" bson:"style,omitempty"`
	MarkerEnd    map[string]any `json:"markerEnd,omitempty" bson:"markerEnd,omitempty"`
}

type Node struct {
	ID               string         `json:"id" bson:"id"`
	Type             string         `json:"type" bson:"type"`
	Position         PositionClass  `json:"position" bson:"position"`
	Data             map[string]any `json:"data" bson:"data"`
	PositionAbsolute PositionClass  `json:"positionAbsolute" bson:"positionAbsolute"`
	Selected         bool           `json:"selected" bson:"selected"`
	Dragging         bool           `json:"dragging" bson:"dragging"`
	Style            map[string]any `json:"style,omitempty" bson:"style,omitempty"`
	ClassName        string         `json:"className" bson:"className"`
	Resizing         *bool          `json:"resizing,omitempty" bson:"resizing,omitempty"`
	ParentNode       *string        `json:"parentNode,omitempty" bson:"parentNode,omitempty"`
	Width            int64          `json:"width" bson:"width"`
	Height           int64          `json:"height" bson:"height"`
}

type PositionClass struct {
	X float64 `json:"x" bson:"x"`
	Y float64 `json:"y" bson:"y"`
}

type Viewport struct {
	X    float64 `json:"x" bson:"x"`
	Y    float64 `json:"y" bson:"y"`
	Zoom float64 `json:"zoom" bson:"zoom"`
}
