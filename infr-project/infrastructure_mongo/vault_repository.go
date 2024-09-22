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

func (r *VaultRepository) InsertMultipleVaultAsync(ctx context.Context, vaults []vault.Vault) <-chan error {
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

func (r *VaultRepository) GetVaultByIdsAsync(ctx context.Context, ids []string) (<-chan []vault.Vault, <-chan error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"id": bson.M{"$in": ids}})
	return r.BaseRepository.GetWithConditionAsync(ctx, condition)
}

func (r *VaultRepository) SearchVaults(ctx context.Context, req *vault.VaultSearch, page *int, size *int) (<-chan []vault.Vault, <-chan error) {
	filter := buildSearchFilter(req)
	condition := extensions.NewSearch(page, size, nil, filter)
	return r.BaseRepository.GetWithConditionAsync(ctx, condition)
}

func buildSearchFilter(req *vault.VaultSearch) map[string]interface{} {
	filter := map[string]interface{}{}
	if req == nil {
		return filter
	}

	if req.Key != "" {
		filter["key"] = bson.D{{Key: "$regex", Value: req.Key}, {Key: "$options", Value: "i"}}
	}

	if req.StorageMedia != "" {
		filter["storage_media"] = req.StorageMedia
	}

	if req.VaultType != "" {
		filter["vault_type"] = req.VaultType
	}

	if req.TypeIdentity != "" {
		filter["type_identity"] = req.TypeIdentity
	}

	if len(req.Tags) > 0 {
		filter["tags"] = bson.D{{Key: "$in", Value: req.Tags}}
	}

	return filter
}
