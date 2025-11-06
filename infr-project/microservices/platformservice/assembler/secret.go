package assembler

import (
	"context"

	"github.com/futugyou/platformservice/application/service"
	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/viewmodel"
)

type SecretAssembler struct{}

func (a *SecretAssembler) ToModel(ctx context.Context, vaultService service.VaultService, secrets []viewmodel.Secret) (map[string]domain.Secret, error) {
	secretInfos := make(map[string]domain.Secret)
	if len(secrets) == 0 {
		return secretInfos, nil
	}

	ids := []string{}
	for _, secret := range secrets {
		ids = append(ids, secret.VaultID)
	}
	vaults, err := vaultService.GetVaultsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, secret := range secrets {
		for i := range vaults {
			if vaults[i].Id == secret.VaultID {
				if vaults[i].VaultType == "system" {
					return nil, err
				}

				secretInfos[secret.Key] = domain.Secret{
					Key:            secret.Key,
					Value:          vaults[i].Id,
					VaultKey:       vaults[i].Key,
					VaultMaskValue: vaults[i].MaskValue,
				}
				break
			}
		}
	}

	return secretInfos, nil
}
