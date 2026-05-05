package did

import (
	"encoding/json"
)

type DescriptorMap struct {
	ID         string      `json:"id"`
	Path       string      `json:"path"`
	Format     string      `json:"format"`
	PathNested *PathNested `json:"path_nested,omitempty"`
}

type PathNested struct {
	ID     string `json:"id"`
	Path   string `json:"path"`
	Format string `json:"format"`
}

type InputDescriptor struct {
	ID          string                           `json:"id"`
	Name        string                           `json:"name"`
	Purpose     string                           `json:"purpose"`
	Format      map[string]InputDescriptorFormat `json:"format"`
	Constraints InputDescriptorConstraints       `json:"constraints"`
}

type InputDescriptorConstraints struct {
	Fields []InputDescriptorConstraintsField `json:"fields"`
}

type InputDescriptorConstraintsField struct {
	ID      string   `json:"id"`
	Path    []string `json:"path"`
	Purpose string   `json:"purpose"`
	Name    string   `json:"name"`
	Filter  any      `json:"filter"`
}

type InputDescriptorFormat struct {
	Alg       []string `json:"alg,omitempty"`
	ProofType []string `json:"proof_type,omitempty"`
}

type PresentationSubmission struct {
	ID            string          `json:"id"`
	DefinitionId  string          `json:"definition_id"`
	DescriptorMap []DescriptorMap `json:"descriptor_map"`
}

type VerifiablePresentationDefinition struct {
	ID               string            `json:"id"`
	InputDescriptors []InputDescriptor `json:"input_descriptors"`
}

type VerifiablePresentation struct {
	Proof                json.RawMessage `json:"proof"`
	Id                   string          `json:"id"`
	Holder               string          `json:"holder"`
	Context              []string        `json:"@context"`
	Types                []string        `json:"type"`
	VerifiableCredential []string        `json:"verifiableCredential"`
}

func (v *VerifiablePresentation) ToDic() map[string]any {
	return map[string]any{
		"@context":             v.Context,
		"id":                   v.Id,
		"type":                 v.Types,
		"holder":               v.Holder,
		"verifiableCredential": v.VerifiableCredential,
	}
}
