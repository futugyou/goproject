package openapi

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/swaggest/openapi-go/openapi31"
	"gopkg.in/yaml.v3"
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
			Title:       swagger.Info.Title,
			Description: &swagger.Info.Description,
			Version:     swagger.Info.Version,
		},
		Paths: &openapi31.Paths{
			MapOfPathItemValues: map[string]openapi31.PathItem{},
			MapOfAnything:       map[string]interface{}{},
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
	return &openapi31.Operation{
		Tags:        op.Tags,
		Summary:     &op.Summary,
		Description: &op.Description,
		Deprecated:  &op.Deprecated,
		// TODO: fill all
	}
}

func saveAsJSON(openAPISpec *openapi31.Spec, filename string) error {
	data, err := json.MarshalIndent(openAPISpec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func saveAsYAML(openAPISpec *openapi31.Spec, filename string) error {
	data, err := yaml.Marshal(openAPISpec)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
