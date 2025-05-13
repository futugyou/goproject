package utilities

import (
	"fmt"
	"strings"
)

type FunctionName struct {
	PluginName string
	Name       string
}

func NewFunctionName(name string, pluginName string) FunctionName {
	return FunctionName{
		Name:       name,
		PluginName: pluginName,
	}
}

func ToFullyQualifiedName(functionName string, pluginName string, functionNameSeparator string) string {
	if len(functionNameSeparator) == 0 {
		functionNameSeparator = "-"
	}
	if len(pluginName) == 0 {
		return functionName
	}
	return fmt.Sprintf("%s%s%s", pluginName, functionNameSeparator, functionName)
}

func ParseFunctionName(fullyQualifiedName string, functionNameSeparator string) FunctionName {

	if len(functionNameSeparator) == 0 {
		functionNameSeparator = "-"
	}

	pluginName := ""
	functionName := fullyQualifiedName

	separatorPos := strings.Index(fullyQualifiedName, functionNameSeparator)
	if separatorPos >= 0 {
		pluginName = strings.TrimSpace(fullyQualifiedName[:separatorPos])
		functionName = strings.TrimSpace(fullyQualifiedName[separatorPos+len(functionNameSeparator):])
	}

	return NewFunctionName(functionName, pluginName)
}
