package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/vault"

	"github.com/futugyou/infr-project/extensions"
)

type VaultRepository struct {
	BaseRepository[vault.Vault]
}

func NewVaultRepository(client *mongo.Client, config DBConfig) *VaultRepository {
	return &VaultRepository{
		BaseRepository: *NewBaseRepository[vault.Vault](client, config),
	}
}

func (r *VaultRepository) InsertMultipleVault(ctx context.Context, vaults []vault.Vault) <-chan error {
	errorChan := make(chan error, 1)
	if len(vaults) == 0 {
		errorChan <- fmt.Errorf("not data need insert")
		return errorChan
	}

	go func() {
		defer close(errorChan)

		c := r.Client.Database(r.DBName).Collection(vaults[0].AggregateName())
		documents := make([]interface{}, 0)
		for i := 0; i < len(vaults); i++ {
			documents = append(documents, vaults[i])
		}
		_, err := c.InsertMany(ctx, documents)
		errorChan <- err
	}()

	return errorChan
}

func (r *VaultRepository) GetAllVaultAsync(ctx context.Context, page *int, size *int) (<-chan []vault.Vault, <-chan error) {
	condition := extensions.NewSearch(page, size, nil, nil)
	return r.BaseRepository.GetWithConditionAsync(ctx, condition)

}

func (r *VaultRepository) GetAllVaultByStorageMediaAsync(ctx context.Context, media vault.StorageMedia) (<-chan []vault.Vault, <-chan error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"storage_media": media.String()})
	return r.BaseRepository.GetWithConditionAsync(ctx, condition)
}

func (r *VaultRepository) GetAllVaultByVaultTypeAsync(ctx context.Context, vType vault.VaultType, identities ...string) (<-chan []vault.Vault, <-chan error) {
	var filter map[string]interface{} = map[string]interface{}{"vault_type": vType.String()}
	if len(identities) > 0 {
		filter["type_identity"] = identities[0]
	}
	condition := extensions.NewSearch(nil, nil, nil, filter)
	return r.BaseRepository.GetWithConditionAsync(ctx, condition)
}

func (r *VaultRepository) GetAllVaultByTagsAsync(ctx context.Context, tags []string) (<-chan []vault.Vault, <-chan error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"tags": bson.M{"$in": tags}})
	return r.BaseRepository.GetWithConditionAsync(ctx, condition)
}

func (r *VaultRepository) GetAllVaultByIdsAsync(ctx context.Context, ids []string) (<-chan []vault.Vault, <-chan error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"id": bson.M{"$in": ids}})
	return r.BaseRepository.GetWithConditionAsync(ctx, condition)
}
