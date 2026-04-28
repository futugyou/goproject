package did

import (
	"context"
	"fmt"
	"strings"
)

type IDidResolver interface {
	Resolve(ctx context.Context, did string) (*DidDocument, error)
	GetMethod() string
}

type IDidFactoryResolver interface {
	Resolve(ctx context.Context, did string) (DidDocument, error)
}

type DidFactoryResolver struct {
	resolvers []IDidResolver
}

func (d *DidFactoryResolver) Resolve(ctx context.Context, did string) (*DidDocument, error) {
	decentralizedIdentifier, err := DidExtractor(did)
	if err != nil {
		return nil, err
	}

	for _, resolve := range d.resolvers {
		if resolve.GetMethod() == decentralizedIdentifier.Method {
			return resolve.Resolve(ctx, did)
		}
	}

	return nil, fmt.Errorf("the method %s doesn't exist", decentralizedIdentifier.Method)
}

func DidExtractor(did string) (*DecentralizedIdentifier, error) {
	if did == "" {
		return nil, fmt.Errorf("invalid DID: empty string")
	}

	splitted := strings.Split(did, ":")
	if len(splitted) != 3 {
		return nil, fmt.Errorf("invalid DID: expected format 'did:method:identifier'")
	}
	if splitted[0] != "did" {
		return nil, fmt.Errorf("invalid DID: expected scheme 'did'")
	}
	identifier := splitted[2]
	fragment := ""

	if strings.Contains(identifier, "#") {
		parts := strings.SplitN(identifier, "#", 2)
		identifier = parts[0]
		fragment = parts[1]
	}

	return &DecentralizedIdentifier{
		Scheme:     splitted[0],
		Method:     splitted[1],
		Identifier: identifier,
		Fragment:   fragment,
	}, nil
}
