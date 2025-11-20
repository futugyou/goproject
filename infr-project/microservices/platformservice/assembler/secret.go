package assembler

import (
	"context"

	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/viewmodel"
)

type SecretAssembler struct{}

func (a *SecretAssembler) ToModel(ctx context.Context, secrets []viewmodel.Secret) (map[string]domain.Secret, error) {
	secretInfos := make(map[string]domain.Secret)
	if len(secrets) == 0 {
		return secretInfos, nil
	}

	for _, secret := range secrets {
		secretInfos[secret.Key] = domain.Secret{
			Key:            secret.Key,
			Value:          secret.VaultID,
			VaultKey:       secret.VaultKey,
			VaultMaskValue: secret.MaskValue,
		}
	}

	return secretInfos, nil
}
