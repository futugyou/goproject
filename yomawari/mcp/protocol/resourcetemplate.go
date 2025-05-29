package protocol

import "strings"

type ResourceTemplate struct {
	UriTemplate string       `json:"uriTemplate"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	MimeType    *string      `json:"mimeType"`
	Annotations *Annotations `json:"annotations"`
}

func (r ResourceTemplate) IsTemplated() bool {
	return strings.ContainsAny(r.UriTemplate, "{")
}

func (r ResourceTemplate) AsResource() *Resource {
	if r.IsTemplated() {
		return nil
	}
	return &Resource{
		Uri:         r.UriTemplate,
		Name:        r.Name,
		Description: r.Description,
		MimeType:    r.MimeType,
		Annotations: r.Annotations,
	}
}
