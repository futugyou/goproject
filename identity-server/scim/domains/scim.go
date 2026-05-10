package domains

import "time"

type Realm struct {
	Name  string
	Owner string
}

type SCIMAttributeMapping struct {
	Id                      string
	SourceAttributeId       string
	SourceResourceType      string
	SourceAttributeSelector string
	TargetResourceType      string
	TargetAttributeId       string
	Mode                    Mode
}

func (s SCIMAttributeMapping) IsSelf() bool {
	return s.SourceResourceType == s.TargetResourceType
}

type SCIMSchemaAttribute struct {
	Id                 string
	FullPath           string
	ParentId           string
	SchemaId           string
	Name               string
	Type               SCIMSchemaAttributeTypes
	MultiValued        bool
	Description        string
	CanonicalValues    []string
	CaseExact          bool
	Mutability         SCIMSchemaAttributeMutabilities
	Returned           SCIMSchemaAttributeReturned
	Uniqueness         SCIMSchemaAttributeUniqueness
	ReferenceTypes     []string
	DefaultValueString []string
	DefaultValueInt    []string
}

type SCIMSchemaExtension struct {
	Id       string
	Schema   string
	Required bool
}

type SCIMSchema struct {
	Id               string
	Name             string
	Description      string
	IsRootSchema     bool
	ResourceType     string
	SchemaExtensions []SCIMSchemaExtension
	Attributes       []SCIMSchemaAttribute
	Representations  []SCIMRepresentation
}

type SCIMRepresentation struct {
	Id             string
	ExternalId     string
	ResourceType   string
	Version        string
	DisplayName    string
	RealmName      string
	Created        time.Time
	LastModified   time.Time
	FlatAttributes []SCIMRepresentationAttribute
	Schemas        []SCIMSchema
}

type SCIMRepresentationAttribute struct {
	Id                 string
	AttributeId        string
	ResourceType       string
	ParentAttributeId  string
	SchemaAttributeId  string
	RepresentationId   string
	FullPath           string
	ValueString        string
	ValueBoolean       bool
	ValueInteger       int
	ValueDateTime      time.Time
	ValueReference     string
	ValueDecimal       float32
	ValueBinary        string
	Namespace          string
	ComputedValueIndex string
	SchemaAttribute    SCIMSchemaAttribute
	Representation     SCIMRepresentation
	Children           []SCIMRepresentationAttribute
	CachedChildren     []SCIMRepresentationAttribute
	IsComputed         bool
}
