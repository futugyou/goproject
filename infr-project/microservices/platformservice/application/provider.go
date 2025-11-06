package application

import (
	"context"
	"fmt"
	"log"

	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/provider"
)

func (s *PlatformService) determineProviderStatus(ctx context.Context, res *domain.Platform) bool {
	provider, err := s.getPlatformProvider(ctx, *res)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	user, err := provider.GetUser(ctx)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	if user == nil || len(user.ID) == 0 {
		log.Printf("no user found for %s provider\n", res.Provider.String())
		return false
	}

	return true
}

func (s *PlatformService) getPlatformProvider(ctx context.Context, src domain.Platform) (provider.PlatformProvider, error) {
	vaultId, err := src.ProviderVaultInfo()
	if err != nil {
		return nil, err
	}

	token, err := s.vaultService.ShowVaultRawValue(ctx, vaultId)
	if err != nil {
		return nil, fmt.Errorf("get platform provider token error, vaultId is %s, message %s", vaultId, err.Error())
	}

	return provider.PlatformProviderFactory(ctx, src.Provider.String(), token)
}
