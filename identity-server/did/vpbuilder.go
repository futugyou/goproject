package did

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
)

type VpBuilderResult struct {
	Vp                     *VerifiablePresentation
	PresentationSubmission *PresentationSubmission
	ErrorMessage           string
}

func (v *VpBuilderResult) HasError() bool {
	return v.ErrorMessage != ""
}

func OKVpBuilderResult(vp *VerifiablePresentation, presentationSubmission *PresentationSubmission) *VpBuilderResult {
	return &VpBuilderResult{
		Vp:                     vp,
		PresentationSubmission: presentationSubmission,
		ErrorMessage:           "",
	}
}

func NokVpBuilderResult(errorMessage string) *VpBuilderResult {
	return &VpBuilderResult{
		Vp:                     nil,
		PresentationSubmission: nil,
		ErrorMessage:           errorMessage,
	}
}

type VcRecord struct {
	Vc             any
	DeserializedVc *W3CVerifiableCredential
	Format         string
	JsonHeader     string
	JsonPayload    string
}

type VpBuilder struct {
	verifiablePresentation *VerifiablePresentation
}

func NewVpBuilder(id string, holder string, vpType *string, jsonLdContext *string) *VpBuilder {
	record := &VerifiablePresentation{
		Id:     id,
		Holder: holder,
	}

	record.Context = append(record.Context, "https://www.w3.org/2018/credentials/v1")
	if jsonLdContext != nil {
		record.Context = append(record.Context, *jsonLdContext)
	}

	record.Context = append(record.Context, "VerifiablePresentation")
	if vpType != nil {
		record.Types = append(record.Types, *vpType)
	}

	return &VpBuilder{
		verifiablePresentation: record,
	}
}

func (b *VpBuilder) BuildAndVerify(vpDef *VerifiablePresentationDefinition, vcRecords []VcRecord, vpFormat string) (*VpBuilderResult, error) {
	var descriptorMaps []DescriptorMap

	if vpDef.InputDescriptors != nil {
		for i, inputDescriptor := range vpDef.InputDescriptors {
			vc := b.getVc(inputDescriptor, vcRecords)
			if vc == nil {
				return nil, fmt.Errorf("VC cannot be resolved")
			}

			if !b.isFormatValid(inputDescriptor, vc) {
				return nil, fmt.Errorf("format not satisfied")
			}

			descriptorMaps = append(descriptorMaps, DescriptorMap{
				ID:     inputDescriptor.ID,
				Path:   "$",
				Format: vpFormat,
				PathNested: &PathNested{
					ID:     inputDescriptor.ID,
					Format: vc.Format,
					Path:   fmt.Sprintf("$.vp.verifiableCredential[%d]", i),
				},
			})
		}
	}

	for _, vc := range vcRecords {
		vcData, _ := json.Marshal(vc.Vc)
		b.verifiablePresentation.VerifiableCredential = append(b.verifiablePresentation.VerifiableCredential, string(vcData))
	}

	return &VpBuilderResult{
		Vp: b.verifiablePresentation,
		PresentationSubmission: &PresentationSubmission{
			ID:            uuid.NewString(),
			DefinitionId:  vpDef.ID,
			DescriptorMap: descriptorMaps,
		},
	}, nil
}

func (b *VpBuilder) isFormatValid(inputDescriptor InputDescriptor, vcRecord *VcRecord) bool {
	if len(inputDescriptor.Format) > 0 {
		for fmtKey, fmtVal := range inputDescriptor.Format {
			if fmtKey != vcRecord.Format {
				continue
			}
			if fmtVal.Alg == nil {
				return true
			}

			vcRecordAlg := gjson.Get(vcRecord.JsonHeader, "alg").String()

			if vcRecordAlg != "" {
				for _, a := range fmtVal.Alg {
					if a == vcRecordAlg {
						return true
					}
				}
			}
		}
		return false
	}
	return true
}

func (b *VpBuilder) getVc(inputDescriptor InputDescriptor, vcRecords []VcRecord) *VcRecord {
	constraints := inputDescriptor.Constraints
	if len(constraints.Fields) == 0 {
		return nil
	}

	for _, vcRecord := range vcRecords {
		isCorrect := true
		for _, field := range constraints.Fields {
			fieldMatch := false
			for _, path := range field.Path {
				res := gjson.Get(vcRecord.JsonPayload, path)

				if res.Exists() {
					if field.Filter != nil {
						if !validateJSONSchema(field.Filter, res.Raw) {
							continue
						}
					}
					fieldMatch = true
					break
				}
			}
			if !fieldMatch {
				isCorrect = false
				break
			}
		}

		if isCorrect {
			return &vcRecord
		}
	}
	return nil
}

func validateJSONSchema(filter any, data string) bool {
	schemaLoader := gojsonschema.NewGoLoader(filter)
	documentLoader := gojsonschema.NewStringLoader(data)
	result, _ := gojsonschema.Validate(schemaLoader, documentLoader)
	return result.Valid()
}
