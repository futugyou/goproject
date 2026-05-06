package did

import (
	"context"
	"fmt"
	"slices"
)

var _ IDidResolver = (*DidKeyResolver)(nil)

type DidKeyResolver struct {
	serializer                     IMulticodecSerializer
	options                        DidKeyOptions
	verificationMethodsStandardLst []IVerificationMethodStandard
}

func NewDidKeyResolver(serializer IMulticodecSerializer, options DidKeyOptions, verificationMethods []IVerificationMethodStandard) *DidKeyResolver {
	return &DidKeyResolver{
		serializer:                     serializer,
		options:                        options,
		verificationMethodsStandardLst: verificationMethods,
	}
}

func NewSimpleDidKeyResolver(options *DidKeyOptions) *DidKeyResolver {
	o := DidKeyOptions{}
	if options != nil {
		o = *options
	}
	return NewDidKeyResolver(FullMulticodecSerializer(), o, GetAllVerificationMethodStandards())
}

// GetMethod implements [IDidResolver].
func (d *DidKeyResolver) GetMethod() string {
	return "key"
}

// Resolve implements [IDidResolver].
func (d *DidKeyResolver) Resolve(ctx context.Context, did string) (*DidDocument, error) {
	decentralizedIdentifier, err := DidExtractor(did)
	if err != nil {
		return nil, err
	}

	if decentralizedIdentifier.Method != "key" {
		return nil, fmt.Errorf("method must be equals to 'key'")
	}

	did = decentralizedIdentifier.GetDidWithoutFragment()
	multibaseValue := decentralizedIdentifier.Identifier
	verificationMethod, err := d.serializer.Deserialize(multibaseValue, "")
	if err != nil {
		return nil, err
	}
	builder := NewDidDocumentBuilderWithID(did)
	verificationMethodId := fmt.Sprint("{0}#{1}", did, multibaseValue)
	publicKeyFormat := d.options.PublicKeyFormat

	if len(publicKeyFormat) == 0 {
		_, ok := verificationMethod.(*JsonWebKeySecurityKey)
		if ok {
			publicKeyFormat = "JsonWebKey2020"
		} else {
			found := false
			for _, v := range d.verificationMethodsStandardLst {
				if slices.Contains(v.GetSupportedCurves(), verificationMethod.GetCrvOrSize()) {
					publicKeyFormat = v.GetType()
					found = true
					break
				}
			}
			if !found {
				publicKeyFormat = "Ed25519VerificationKey2020"
			}
		}
	}

	builder.AddVerificationMethod(publicKeyFormat,
		verificationMethod,
		did,
		Authentication|AssertionMethod|CapabilityInvocation|CapabilityDelegation,
		nil,
		nil,
		nil,
		func(c *DidDocumentVerificationMethod) {
			c.ID = verificationMethodId
		})

	if d.options.EnableEncryptionKeyDerivation {
		return nil, fmt.Errorf("This feature is not yet supported")

	}

	return builder.Build(), nil
}
