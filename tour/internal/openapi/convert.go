package openapi

import (
	"os"
	"strconv"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/swaggest/openapi-go/openapi31"

	util "github.com/futugyou/extensions"
)

func ConvertSwaggerToOpenapi(path string, outpath string) error {
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

	if strings.HasSuffix(outpath, ".yaml") {
		err = saveAsYAML(openAPISpec, outpath)
	}
	if strings.HasSuffix(outpath, ".json") {
		err = saveAsJSON(openAPISpec, outpath)
	}
	if err != nil {
		return err
	}
	return nil
}

func convertSwaggerToOpenAPI(swagger *spec.Swagger) *openapi31.Spec {
	openAPI := &openapi31.Spec{
		Openapi: "3.0.0",
		Info: openapi31.Info{
			Title:          swagger.Info.Title,
			Description:    &swagger.Info.Description,
			TermsOfService: &swagger.Info.TermsOfService,
			Contact: &openapi31.Contact{
				Name:  util.GetStringFieldPointer(swagger, "Info", "Contact", "Name"),
				URL:   util.GetStringFieldPointer(swagger, "Info", "Contact", "URL"),
				Email: util.GetStringFieldPointer(swagger, "Info", "Contact", "Email"),
			},
			License: &openapi31.License{
				Name: util.GetStringFieldStruct(swagger, "Info", "License", "Name"),
				URL:  util.GetStringFieldPointer(swagger, "Info", "License", "URL"),
			},
			Version: swagger.Info.Version,
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

	return openAPI
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
		ID:           &op.ID,
		Parameters:   []openapi31.ParameterOrReference{},
		RequestBody:  &openapi31.RequestBodyOrReference{},
		Responses:    convertResponses(op.Responses, op.Produces),
		Deprecated:   &op.Deprecated,
		Security:     op.Security,
	}
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
				Content:     convertContent(produces, v.ResponseProps),
				Links:       map[string]openapi31.LinkOrReference{},
			},
		}
	}
	return &openapi31.Responses{
		MapOfResponseOrReferenceValues: maps,
		MapOfAnything:                  responses.Extensions,
	}
}

func convertContent(produces []string, responseProps spec.ResponseProps) map[string]openapi31.MediaType {
	mediaTypes := make(map[string]openapi31.MediaType)
	for _, v := range produces {
		if responseProps.Schema == nil {
			continue
		}

		schema := make(map[string]interface{})
		ref := responseProps.Schema.Ref.Ref.String()
		if len(ref) > 0 {
			schema["$ref"] = ref
		} else {
			schema["type"] = responseProps.Schema.Type
			schema["items"] = responseProps.Schema.Items
		}

		mediaTypes[v] = openapi31.MediaType{
			Schema: schema,
		}
	}

	return mediaTypes
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
	return os.WriteFile(filename, data, 0644)
}

func saveAsYAML(openAPISpec *openapi31.Spec, filename string) error {
	data, err := openAPISpec.MarshalYAML()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
