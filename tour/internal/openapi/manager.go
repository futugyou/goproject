package openapi

import (
	"os"
	"reflect"
	"strings"

	"github/go-project/tour/util"

	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi31"
)

type Manager struct {
	Infos  []OpenAPIOperation
	Config OpenAPIConfig
}

func NewManager(astManager util.ASTManager, config OpenAPIConfig) (*Manager, error) {
	sList := make([]OpenAPIOperation, 0)
	for _, api := range config.APIConfigs {
		req, err := astManager.GetReflectTypeByName(api.Request)
		if err != nil {
			return nil, err
		}

		if api.Method == "DELETE" || api.Method == "GET" {
			req = astManager.ConvertReflectTypeTag(req, "json", "query")
		}

		resp, err := astManager.GetReflectTypeByName(api.Response)
		if err != nil {
			return nil, err
		}

		o := NewOpenAPIOperation(api.Method, api.Path, api.Description, req, resp)
		sList = append(sList, *o)
	}

	m := &Manager{
		Infos:  sList,
		Config: config,
	}

	return m, nil
}

func (m *Manager) GenerateOpenAPI() error {
	reflector := openapi31.NewReflector()
	reflector.Spec = &openapi31.Spec{Openapi: m.Config.SpceVersion}
	reflector.Spec.Info.
		WithTitle(m.Config.Title).
		WithVersion(m.Config.APIVersion).
		WithDescription(m.Config.Description)

	reflector.Reflector.DefaultOptions = append(reflector.Reflector.DefaultOptions,
		jsonschema.InterceptNullability(func(params jsonschema.InterceptNullabilityParams) {
			// Removing nullability from non-pointer slices (regardless of omitempty).
			if params.Type.Kind() != reflect.Ptr && params.Schema.HasType(jsonschema.Null) && params.Schema.HasType(jsonschema.Array) {
				*params.Schema.Type = jsonschema.Array.Type()
			}
		}))

	for _, info := range m.Infos {
		op, err := reflector.NewOperationContext(info.Method, info.Path)
		if err != nil {
			return err
		}
		op.SetDescription(info.Description)
		op.AddReqStructure(info.Request)
		op.AddRespStructure(info.Response)
		if err = reflector.AddOperation(op); err != nil {
			return err
		}
	}

	return m.dumpOpenAPISpec(reflector.Spec)
}

func (m *Manager) dumpOpenAPISpec(spec *openapi31.Spec) error {
	var schema []byte
	var err error
	ot := strings.ToLower(m.Config.OutputType)
	if ot == "yaml" || ot == "yml" {
		schema, err = spec.MarshalYAML()
	} else {
		schema, err = spec.MarshalJSON()
	}

	if err != nil {
		return err
	}

	return os.WriteFile(m.Config.OutputPath, schema, 0600)
}
