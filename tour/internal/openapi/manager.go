package openapi

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github/go-project/tour/util"

	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
)

type Manager struct {
	Infos  []OpenAPIOperation
	Config OpenAPIConfig
}

func NewManager(ts []util.StructInfo, config OpenAPIConfig) (*Manager, error) {
	sList := make([]OpenAPIOperation, 0)
	for _, js := range config.JsonConfig {
		req, err := util.GetReflectTypeFromStructInfo(js.Request, ts)
		if err != nil {
			return nil, err
		}
		resp, err := util.GetReflectTypeFromStructInfo(js.Response, ts)
		if err != nil {
			return nil, err
		}
		o := NewOpenAPIOperation(js.Method, js.Path, js.Description, req, resp)
		sList = append(sList, *o)
	}
	m := &Manager{
		Infos:  sList,
		Config: config,
	}
	return m, nil
}

func (m *Manager) GenerateOpenAPI() error {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: m.Config.SpceVersion}
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

func (m *Manager) dumpOpenAPISpec(spec *openapi3.Spec) error {
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
	fmt.Println(m.Config.OutputPath)
	return os.WriteFile(m.Config.OutputPath, schema, 0600)
}
