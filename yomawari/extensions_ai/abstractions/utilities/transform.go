package utilities

import (
	"encoding/json"
	"fmt"
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
				return map[string]interface{}{"not": true}
			}
			return map[string]interface{}{}
		}
		return schema

	case map[string]interface{}:
		obj := val

		var properties map[string]interface{}
		if rawProps, ok := obj["properties"].(map[string]interface{}); ok {
			properties = rawProps
			path = append(path, "properties")
			for key, v := range rawProps {
				path = append(path, key)
				rawProps[key] = transformSchemaCore(v, options, path)
				path = path[:len(path)-1]
			}
			path = path[:len(path)-1]
		}

		if items, ok := obj["items"]; ok {
			path = append(path, "items")
			obj["items"] = transformSchemaCore(items, options, path)
			path = path[:len(path)-1]
		}

		if addProps, ok := obj["additionalProperties"]; ok {
			if b, isBool := addProps.(bool); !isBool || b {
				path = append(path, "additionalProperties")
				obj["additionalProperties"] = transformSchemaCore(addProps, options, path)
				path = path[:len(path)-1]
			}
		}

		if notSchema, ok := obj["not"]; ok {
			path = append(path, "not")
			obj["not"] = transformSchemaCore(notSchema, options, path)
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
			if _, exists := obj["additionalProperties"]; !exists {
				obj["additionalProperties"] = false
			}
		}

		if options.RequireAllProperties && properties != nil {
			required := make([]interface{}, 0, len(properties))
			for k := range properties {
				required = append(required, k)
			}
			obj["required"] = required
		}

		if options.UseNullableKeyword {
			if typeArray, ok := obj["type"].([]interface{}); ok {
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
					obj["type"] = foundType
					obj["nullable"] = true
				}
			}
		}

		if options.MoveDefaultKeywordToDescription {
			if def, ok := obj["default"]; ok {
				desc := ""
				if d, ok := obj["description"].(string); ok {
					desc = d + " "
				}
				defJson, _ := json.Marshal(def)
				desc += fmt.Sprintf("(Default value: %s)", string(defJson))
				obj["description"] = strings.TrimSpace(desc)
				delete(obj, "default")
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
