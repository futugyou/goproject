package ollama

import "encoding/json"

type OllamaFunctionToolParameter struct {
	Description *string
	Type        *string
	Enum        []string
}

//
// func main() {
//     jsonStr := `{"Properties":{"key":"value"}}`
//     var parameters OllamaFunctionToolParameters
//     json.Unmarshal([]byte(jsonStr), &parameters)

//	    var rawMap map[string]interface{}
//	    json.Unmarshal(parameters.Properties, &rawMap)
//	    println(rawMap["key"].(string)) // output: value
//	}
type OllamaFunctionToolParameters struct {
	Type       string
	Required   []string
	Properties map[string]json.RawMessage
}
