package utilities

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func TransformSchema(schema json.RawMessage, options AIJsonSchemaTransformOptions) (json.RawMessage, error) {
	var schemaMap interface{}
	if err := json.Unmarshal(schema, &schemaMap); err != nil {
		return nil, err
	}

	transformed := transformSchemaCore(schemaMap, options, nil)

	return json.Marshal(transformed)
}

func transformSchemaCore(schema interface{}, options AIJsonSchemaTransformOptions, path []string) interface{} {
	switch val := schema.(type) {
	case bool:
		if options.ConvertBooleanSchemas {
			if !val {
				return map[string]interface{}{NotPropertyName: true}
			}
			return map[string]interface{}{}
		}
		return schema

	case map[string]interface{}:
		obj := val

		var properties map[string]interface{}
		if rawProps, ok := obj[PropertiesPropertyName].(map[string]interface{}); ok {
			properties = rawProps
			path = append(path, PropertiesPropertyName)
			for key, v := range rawProps {
				path = append(path, key)
				rawProps[key] = transformSchemaCore(v, options, path)
				path = path[:len(path)-1]
			}
			path = path[:len(path)-1]
		}

		if items, ok := obj[ItemsPropertyName]; ok {
			path = append(path, ItemsPropertyName)
			obj[ItemsPropertyName] = transformSchemaCore(items, options, path)
			path = path[:len(path)-1]
		}

		if addProps, ok := obj[AdditionalPropertiesName]; ok {
			if b, isBool := addProps.(bool); !isBool || b {
				path = append(path, AdditionalPropertiesName)
				obj[AdditionalPropertiesName] = transformSchemaCore(addProps, options, path)
				path = path[:len(path)-1]
			}
		}

		if notSchema, ok := obj[NotPropertyName]; ok {
			path = append(path, NotPropertyName)
			obj[NotPropertyName] = transformSchemaCore(notSchema, options, path)
			path = path[:len(path)-1]
		}

		for _, keyword := range []string{"anyOf", "oneOf", "allOf"} {
			if arr, ok := obj[keyword].([]interface{}); ok {
				path = append(path, keyword)
				for i, item := range arr {
					path = append(path, fmt.Sprintf("[%d]", i))
					arr[i] = transformSchemaCore(item, options, path)
					path = path[:len(path)-1]
				}
				path = path[:len(path)-1]
			}
		}

		if options.DisallowAdditionalProperties && properties != nil {
			if _, exists := obj[AdditionalPropertiesName]; !exists {
				obj[AdditionalPropertiesName] = false
			}
		}

		if options.RequireAllProperties && properties != nil {
			required := make([]interface{}, 0, len(properties))
			for k := range properties {
				required = append(required, k)
			}
			obj[RequiredPropertyName] = required
		}

		if options.UseNullableKeyword {
			if typeArray, ok := obj[TypePropertyName].([]interface{}); ok {
				var foundType string
				isNullable := false

				for _, t := range typeArray {
					if str, ok := t.(string); ok {
						if str == "null" {
							isNullable = true
						} else if foundType == "" {
							foundType = str
						} else {
							foundType = ""
							break
						}
					}
				}

				if isNullable && foundType != "" {
					obj[TypePropertyName] = foundType
					obj["nullable"] = true
				}
			}
		}

		if options.MoveDefaultKeywordToDescription {
			if def, ok := obj[DefaultPropertyName]; ok {
				desc := ""
				if d, ok := obj[DescriptionPropertyName].(string); ok {
					desc = d + " "
				}
				defJson, _ := json.Marshal(def)
				desc += fmt.Sprintf("(Default value: %s)", string(defJson))
				obj[DescriptionPropertyName] = strings.TrimSpace(desc)
				delete(obj, DefaultPropertyName)
			}
		}

		if options.TransformSchemaNode != nil {
			schema = options.TransformSchemaNode(AIJsonSchemaTransformContext{path: append([]string{}, path...)}, obj)
		} else {
			schema = obj
		}

	default:
		panic("schema must be a boolean or object")
	}

	return schema
}

func TransformSchemaNode(ctx AIJsonSchemaCreateContext, schema map[string]interface{}) map[string]interface{} {
	localDescription := ""
	if ctx.Path.IsEmpty() && ctx.Description != "" {
		localDescription = ctx.Description
	} else if descAttr := ctx.GetCustomDescription(); descAttr != "" {
		localDescription = descAttr
	}

	if ctx.ParameterName != "" {
		if refVal, ok := schema[RefPropertyName].(string); ok {
			if refVal == "#" {
				schema[RefPropertyName] = fmt.Sprintf("#/properties/%s", ctx.ParameterName)
			} else if strings.HasPrefix(refVal, "#/") {
				schema[RefPropertyName] = fmt.Sprintf("#/properties/%s/%s", ctx.ParameterName, refVal[len("#/"):])
			}
		}
	}

	if isEnum(ctx.Type) {
		if _, hasEnum := schema[EnumPropertyName]; hasEnum {
			if _, hasType := schema[TypePropertyName]; !hasType {
				schema = insertAtStart(schema, TypePropertyName, "string")
			}
		}
	}

	if isNullableEnum(ctx.Type) {
		if _, hasEnum := schema[EnumPropertyName]; hasEnum {
			if _, hasType := schema[TypePropertyName]; !hasType {
				schema = insertAtStart(schema, TypePropertyName, []interface{}{"string", "null"})
			}
		}
	}

	for _, keyword := range schemaKeywordsDisallowedByVendors {
		delete(schema, keyword)
	}

	if t, ok := typeIsIntegerWithStringNumberHandling(ctx, schema); ok {
		schema[TypePropertyName] = t
		delete(schema, PatternPropertyName)
	}

	if ctx.Path.IsEmpty() && ctx.HasDefaultValue {
		schema[DefaultPropertyName] = ctx.SerializeDefaultValue()
	}

	if localDescription != "" {
		schema = insertAtStart(schema, DescriptionPropertyName, localDescription)
	}

	if ctx.Path.IsEmpty() && ctx.InferenceOptions.IncludeSchemaKeyword {
		schema = insertAtStart(schema, SchemaPropertyName, SchemaKeywordUri)
	}

	if ctx.InferenceOptions.TransformSchemaNode != nil {
		schema = ctx.InferenceOptions.TransformSchemaNode(ctx, schema)
	}

	return schema
}
func insertAtStart(obj map[string]interface{}, key string, value interface{}) map[string]interface{} {
	// Since the map is unordered, the insertion order cannot be controlled and can usually only be processed during the serialization phase
	obj[key] = value
	return obj
}

func isEnum(typ reflect.Type) bool {
	return typ.Kind() == reflect.Int && typ.Name() != ""
}

func isNullableEnum(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		return isEnum(typ.Elem())
	}
	return false
}

func typeIsIntegerWithStringNumberHandling(ctx AIJsonSchemaCreateContext, schema map[string]interface{}) (string, bool) {
	typ, ok := schema[TypePropertyName]
	arr, isArr := typ.([]interface{})
	if !ok || !isArr {
		return "", false
	}

	hasString := false
	hasNumber := false
	for _, v := range arr {
		if str, ok := v.(string); ok {
			switch str {
			case "string":
				hasString = true
			case "number", "integer":
				hasNumber = true
			}
		}
	}
	if hasNumber && hasString {
		if strings.Contains(ctx.Type.String(), "int") {
			return "integer", true
		}
		return "number", true
	}
	return "", false
}
