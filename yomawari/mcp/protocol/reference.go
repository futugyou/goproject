package protocol

import "fmt"

type Reference struct {
	Type string  `json:"type"`
	Url  *string `json:"url"`
	Name *string `json:"name"`
}

func (r *Reference) ToString() string {
	name := ""
	if r.Url != nil {
		name = *r.Url
	} else if r.Name != nil {
		name = *r.Name
	}
	return fmt.Sprintf("\"%s\": \"%s\"", r.Type, name)
}

func (r *Reference) Validate() (string, bool) {
	validationMessage := ""
	if r.Type == "ref/resource" {
		if r.Url == nil || len(*r.Url) == 0 {
			validationMessage = "Uri is required for ref/resource"
			return validationMessage, false
		}
	} else if r.Type == "ref/prompt" {
		if r.Name != nil || len(*r.Name) == 0 {
			validationMessage = "Name is required for ref/prompt"
			return validationMessage, false
		}
	} else {
		validationMessage = "Unknown reference type: " + r.Type
		return validationMessage, false
	}

	return validationMessage, true
}
