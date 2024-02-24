package mongo2struct

import (
	"github/go-project/tour/internal/word"
	"slices"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
)

func UnderscoreToUpperCamelCase(s string) string {
	s = word.UnderscoreToUpperCamelCase(s)
	if unicode.IsDigit(rune(s[0])) {
		s = "A" + s
	}
	return s
}

func createImports(s []string, itemType string) []string {
	// for now 'time' only
	if itemType == "time.Time" && !slices.Contains(s, "\"time\"") {
		s = append(s, "\"time\"")
	}

	return s
}

func convertBsontypeTogotype(value bson.RawValue) string {
	switch value.Type {
	case bson.TypeDouble:
		return "float64"
	case bson.TypeString:
		return "string"
	case bson.TypeEmbeddedDocument:
	case bson.TypeArray:
		e, err := value.Array().Elements()
		if err != nil || len(e) == 0 {
			return "[]interface{}"
		}
		return "[]" + convertBsontypeTogotype(e[0].Value())
	case bson.TypeBinary:
	case bson.TypeUndefined:
	case bson.TypeObjectID:
		return "string"
	case bson.TypeBoolean:
		return "bool"
	case bson.TypeDateTime:
		return "time.Time"
	case bson.TypeNull:
	case bson.TypeRegex:
	case bson.TypeDBPointer:
	case bson.TypeJavaScript:
	case bson.TypeSymbol:
	case bson.TypeCodeWithScope:
	case bson.TypeInt32:
		return "int32"
	case bson.TypeTimestamp:
		return "int64"
	case bson.TypeInt64:
		return "int64"
	case bson.TypeDecimal128:
		return "float64"
	case bson.TypeMinKey:
	case bson.TypeMaxKey:
	}

	return "interface{}"
}
