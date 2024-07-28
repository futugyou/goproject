package openapi

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/swaggest/openapi-go/openapi31"

	util "github.com/futugyou/extensions"
)

func ConvertSwaggerToOpenapi(path string, outpath string, fileType string) error {
	swaggerData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	doc, err := loads.Analyzed(swaggerData, "")
	if err != nil {
		return err
	}

	swagger := doc.Spec()
	openAPISpec := convertSwaggerToOpenAPI(swagger)

	if fileType == "yaml" {
		path := adjustFileType(outpath, fileType)
		err = saveAsYAML(openAPISpec, path)
	} else if fileType == "json" {
		path := adjustFileType(outpath, fileType)
		err = saveAsJSON(openAPISpec, path)
	} else {
		if strings.HasSuffix(outpath, ".yaml") {
			err = saveAsYAML(openAPISpec, outpath)
		}
		if strings.HasSuffix(outpath, ".json") {
			err = saveAsJSON(openAPISpec, outpath)
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func adjustFileType(outpath, fileType string) string {
	ext := filepath.Ext(outpath)
	return strings.TrimSuffix(outpath, ext) + "." + fileType
}

func convertSwaggerToOpenAPI(swagger *spec.Swagger) *openapi31.Spec {
	openAPI := &openapi31.Spec{
		Openapi: "3.0.0",
		Info: openapi31.Info{
			Title:          util.GetStringFieldStruct(swagger, "Info", "Title"),
			Description:    util.GetStringFieldPointer(swagger, "Info", "Description"),
			TermsOfService: util.GetStringFieldPointer(swagger, "Info", "TermsOfService"),
			Contact: &openapi31.Contact{
				Name:  util.GetStringFieldPointer(swagger, "Info", "Contact", "Name"),
				URL:   util.GetStringFieldPointer(swagger, "Info", "Contact", "URL"),
				Email: util.GetStringFieldPointer(swagger, "Info", "Contact", "Email"),
			},
			License: &openapi31.License{
				Name: util.GetStringFieldStruct(swagger, "Info", "License", "Name"),
				URL:  util.GetStringFieldPointer(swagger, "Info", "License", "URL"),
			},
			Version: util.GetStringFieldStruct(swagger, "Info", "Version"),
		},
		Servers: []openapi31.Server{
			{
				URL: "/",
			},
		},
		Security: swagger.Security,
		Paths: &openapi31.Paths{
			MapOfPathItemValues: map[string]openapi31.PathItem{},
		},
	}

	for path, pathItem := range swagger.Paths.Paths {
		openAPI.Paths.MapOfPathItemValues[path] = openapi31.PathItem{
			Get:    convertOperation(pathItem.Get),
			Post:   convertOperation(pathItem.Post),
			Put:    convertOperation(pathItem.Put),
			Delete: convertOperation(pathItem.Delete),
			// TODO: add more
		}
	}

	openAPI.Components = convertComponents(swagger.Definitions)

	return openAPI
}

func convertComponents(defs spec.Definitions) *openapi31.Components {
	coms := &openapi31.Components{
		Schemas:         map[string]map[string]interface{}{},
		Responses:       map[string]openapi31.ResponseOrReference{},
		Parameters:      map[string]openapi31.ParameterOrReference{},
		Examples:        map[string]openapi31.ExampleOrReference{},
		RequestBodies:   map[string]openapi31.RequestBodyOrReference{},
		Headers:         map[string]openapi31.HeaderOrReference{},
		SecuritySchemes: map[string]openapi31.SecuritySchemeOrReference{},
		Links:           map[string]openapi31.LinkOrReference{},
		Callbacks:       map[string]openapi31.CallbacksOrReference{},
		PathItems:       map[string]openapi31.PathItemOrReference{},
	}

	for k, v := range defs {
		coms.Schemas[k] = make(map[string]interface{})
		if len(v.Required) > 0 {
			coms.Schemas[k]["required"] = v.Required
		}
		coms.Schemas[k]["type"] = v.Type
		coms.Schemas[k]["properties"] = v.Properties
	}

	return coms
}

func convertOperation(op *spec.Operation) *openapi31.Operation {
	if op == nil {
		return nil
	}
	description := util.GetStringFieldPointer(op, "ExternalDocs", "Description")
	url := util.GetStringFieldStruct(op, "ExternalDocs", "URL")
	var externalDocs *openapi31.ExternalDocumentation = nil
	if (description != nil && len(*description) > 0) || (len(url) > 0) {
		externalDocs = &openapi31.ExternalDocumentation{
			Description: description,
			URL:         url,
		}
	}

	return &openapi31.Operation{
		Tags:         op.Tags,
		Summary:      &op.Summary,
		Description:  &op.Description,
		ExternalDocs: externalDocs,
		ID:           nil,
		Parameters:   convertParameter(op.Parameters),
		RequestBody:  convertRequestBody(op.Parameters, op.Produces),
		Responses:    convertResponses(op.Responses, op.Produces),
		Deprecated:   &op.Deprecated,
		Security:     op.Security,
	}
}

func convertRequestBody(parameters []spec.Parameter, produces []string) *openapi31.RequestBodyOrReference {
	if len(parameters) == 0 {
		return nil
	}

	var parameter *spec.Parameter = nil
	for _, p := range parameters {
		if p.In == "body" {
			parameter = &p
			break
		}
	}
	if parameter == nil {
		return nil
	}
	result := &openapi31.RequestBodyOrReference{
		RequestBody: &openapi31.RequestBody{
			Description: &parameter.Description,
			Content:     convertContent(produces, parameter.Schema, parameter.SimpleSchema),
			Required:    &parameter.Required,
		},
	}

	return result
}

func convertParameter(parameters []spec.Parameter) []openapi31.ParameterOrReference {
	if len(parameters) == 0 {
		return nil
	}
	result := make([]openapi31.ParameterOrReference, 0)
	for _, parameter := range parameters {
		if parameter.In == "body" {
			continue
		}

		p := openapi31.ParameterOrReference{
			Parameter: &openapi31.Parameter{
				Name:        parameter.Name,
				In:          openapi31.ParameterIn(parameter.In),
				Description: &parameter.Description,
				Required:    &parameter.Required,
				Schema:      convertSchema(parameter.Schema, parameter.SimpleSchema),
				Example:     convertExample(parameter.Example),
			},
		}
		result = append(result, p)
	}
	return result
}

func convertExample(example interface{}) *interface{} {
	if example == nil {
		return nil
	}
	return &example
}

func convertResponses(responses *spec.Responses, produces []string) *openapi31.Responses {
	if responses == nil {
		return nil
	}
	maps := make(map[string]openapi31.ResponseOrReference)
	for k, v := range responses.StatusCodeResponses {
		key := strconv.Itoa(k)
		maps[key] = openapi31.ResponseOrReference{
			Response: &openapi31.Response{
				Description: v.Description,
				Headers:     convertHeaders(v.Headers),
				Content:     convertContent(produces, v.ResponseProps.Schema, spec.SimpleSchema{}),
				Links:       map[string]openapi31.LinkOrReference{},
			},
		}
	}
	return &openapi31.Responses{
		MapOfResponseOrReferenceValues: maps,
		MapOfAnything:                  responses.Extensions,
	}
}

func convertContent(produces []string, schema *spec.Schema, simpleSchema spec.SimpleSchema) map[string]openapi31.MediaType {
	if schema == nil {
		if len(simpleSchema.Type) == 0 {
			return nil
		}
	}

	mediaTypes := make(map[string]openapi31.MediaType)
	for _, v := range produces {
		mediaTypes[v] = openapi31.MediaType{
			Schema: convertSchema(schema, simpleSchema),
		}
	}

	return mediaTypes
}

func convertSchema(swaggerSchema *spec.Schema, simpleSchema spec.SimpleSchema) map[string]interface{} {
	if swaggerSchema == nil {
		if len(simpleSchema.Type) == 0 {
			return nil
		}
	}

	schema := make(map[string]interface{})
	if swaggerSchema != nil {
		ref := swaggerSchema.Ref.String()
		if len(ref) > 0 {
			schema["$ref"] = ref
		} else {
			schema["type"] = swaggerSchema.Type
			if swaggerSchema.Items == nil || swaggerSchema.Items.Schema == nil {
				return schema
			}

			ref = swaggerSchema.Items.Schema.Ref.String()
			if len(ref) == 0 {
				return schema
			}

			schema["items"] = map[string]interface{}{
				"$ref": ref,
			}
		}
	} else {
		schema["type"] = simpleSchema.Type
		if simpleSchema.Items != nil {
			schema["items"] = simpleSchema.Items
		}
	}

	return schema
}

func convertRef(ref string) string {
	return strings.ReplaceAll(ref, "#/definitions/", "#/components/schemas/")
}

func convertHeaders(header map[string]spec.Header) map[string]openapi31.HeaderOrReference {
	headers := make(map[string]openapi31.HeaderOrReference)
	for kk, vv := range header {
		headers[kk] = openapi31.HeaderOrReference{
			Reference: &openapi31.Reference{},
			Header: &openapi31.Header{
				Description:   &vv.Description,
				Example:       &vv.Example,
				MapOfAnything: vv.Extensions,
				Schema: map[string]interface{}{
					vv.SimpleSchema.CollectionFormat: vv.SimpleSchema,
				},
			},
		}
	}
	return headers
}

func saveAsJSON(openAPISpec *openapi31.Spec, filename string) error {
	data, err := openAPISpec.MarshalJSON()
	if err != nil {
		return err
	}
	data = []byte(convertRef(string(data)))
	return os.WriteFile(filename, data, 0644)
}

func saveAsYAML(openAPISpec *openapi31.Spec, filename string) error {
	data, err := openAPISpec.MarshalYAML()
	if err != nil {
		return err
	}
	data = []byte(convertRef(string(data)))
	return os.WriteFile(filename, data, 0644)
}
