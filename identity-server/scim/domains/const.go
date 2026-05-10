package domains

const (
	ResourceTypeAttributeSchemas          string = "schemas"
	ResourceTypeAttributeId               string = "id"
	ResourceTypeAttributeName             string = "name"
	ResourceTypeAttributeDescription      string = "description"
	ResourceTypeAttributeEndpoint         string = "endpoint"
	ResourceTypeAttributeSchemaExtensions string = "schemaExtensions"
	ResourceTypeAttributeSchema           string = "schema"
	ResourceTypeAttributeRequired         string = "required"
	ResourceTypeAttributeMeta             string = "meta"
)

const (
	SCIMEndpointsUser                  string = "Users"
	SCIMEndpointsGroup                 string = "Groups"
	SCIMEndpointsServiceProviderConfig string = "ServiceProviderConfig"
	SCIMEndpointsBulk                  string = "Bulk"
	SCIMEndpointsSchemas               string = "Schemas"
	SCIMEndpointsResourceTypes         string = "ResourceTypes"
	SCIMEndpointsProvisioning          string = "Provisioning"
)

const (
	SCIMResourceTypesUser                  string = "Users"
	SCIMResourceTypesGroup                 string = "Group"
	SCIMResourceTypesResourceType          string = "ResourceType"
	SCIMResourceTypesServiceProviderConfig string = "ServiceProviderConfig"
	SCIMResourceTypesSchema                string = "Schema"
)

const (
	StandardSCIMMetaAttributesResourceType string = "resourceType"
	StandardSCIMMetaAttributesCreated      string = "created"
	StandardSCIMMetaAttributesLastModified string = "lastModified"
	StandardSCIMMetaAttributesLocation     string = "location"
	StandardSCIMMetaAttributesVersion      string = "version"
)

const (
	StandardSCIMRepresentationAttributesSchemas         = "schemas"
	StandardSCIMRepresentationAttributesMeta            = "meta"
	StandardSCIMRepresentationAttributesId              = "id"
	StandardSCIMRepresentationAttributesName            = "name"
	StandardSCIMRepresentationAttributesDescription     = "description"
	StandardSCIMRepresentationAttributesAttributes      = "attributes"
	StandardSCIMRepresentationAttributesExternalId      = "externalId"
	StandardSCIMRepresentationAttributesTotalResults    = "totalResults"
	StandardSCIMRepresentationAttributesStartIndex      = "startIndex"
	StandardSCIMRepresentationAttributesItemsPerPage    = "itemsPerPage"
	StandardSCIMRepresentationAttributesResources       = "Resources"
	StandardSCIMRepresentationAttributesOperations      = "Operations"
	StandardSCIMRepresentationAttributesMethod          = "method"
	StandardSCIMRepresentationAttributesPath            = "path"
	StandardSCIMRepresentationAttributesBulkId          = "bulkId"
	StandardSCIMRepresentationAttributesData            = "data"
	StandardSCIMRepresentationAttributesLocation        = "location"
	StandardSCIMRepresentationAttributesVersion         = "version"
	StandardSCIMRepresentationAttributesType            = "type"
	StandardSCIMRepresentationAttributesMultiValued     = "multiValued"
	StandardSCIMRepresentationAttributesRequired        = "required"
	StandardSCIMRepresentationAttributesCaseExact       = "caseExact"
	StandardSCIMRepresentationAttributesMutability      = "mutability"
	StandardSCIMRepresentationAttributesReturned        = "returned"
	StandardSCIMRepresentationAttributesUniqueness      = "uniqueness"
	StandardSCIMRepresentationAttributesSubAttributes   = "subAttributes"
	StandardSCIMRepresentationAttributesCanonicalValues = "canonicalValues"
	StandardSCIMRepresentationAttributesDisplayName     = "displayName"
)
