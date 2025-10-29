package domain

// ResourceType is the interface for resource types.
type ResourceType interface {
	privateResourceType() // Prevents external implementation
	String() string
}

// resourceType is the underlying implementation for ResourceType.
type resourceType string

// privateResourceType makes resourceType implement ResourceType.
func (c resourceType) privateResourceType() {}

// String makes resourceType implement ResourceType.
func (c resourceType) String() string {
	return string(c)
}

// Constants for the different resource types.
const (
	DrawIO     resourceType = "DrawIO"
	Markdown   resourceType = "Markdown"
	Excalidraw resourceType = "Excalidraw"
	Plate      resourceType = "Plate"
)

func GetResourceType(rType string) ResourceType {
	switch rType {
	case "DrawIO":
		return DrawIO
	case "Markdown":
		return Markdown
	case "Excalidraw":
		return Excalidraw
	case "Plate":
		return Plate
	default:
		return Markdown
	}
}
