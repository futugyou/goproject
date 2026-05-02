package did

import (
	"github.com/piprate/json-gold/ld"
)

func RdfTransform(input any) (string, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Format = "application/n-quads"
	options.Algorithm = "URDNA2015"

	normalized, err := proc.Normalize(input, options)
	if err != nil {
		return "", err
	}

	return normalized.(string), nil
}
