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

func (r *VaultRepository) SearchVaults(ctx context.Context, req []vault.VaultSearch, page *int, size *int) (<-chan []vault.Vault, <-chan error) {
	filter := buildSearchFilter(req)
	condition := extensions.NewSearch(page, size, nil, filter)
	return r.BaseRepository.GetWithConditionAsync(ctx, condition)
}

func buildSearchFilter(reqs []vault.VaultSearch) map[string]interface{} {
	filter := map[string]interface{}{}
	if len(reqs) == 0 {
		return filter
	}

	orConditions := bson.A{}
	for _, req := range reqs {
		andConditions := bson.D{}

		if req.ID != "" {
			andConditions = append(andConditions, bson.E{Key: "id", Value: req.ID})
		}

		if req.Key != "" {
			if req.KeyFuzzy {
				andConditions = append(andConditions, bson.E{Key: "key", Value: bson.D{{Key: "$regex", Value: req.Key}, {Key: "$options", Value: "i"}}})
			} else {
				andConditions = append(andConditions, bson.E{Key: "key", Value: req.Key})
			}
		}

		if req.StorageMedia != "" {
			andConditions = append(andConditions, bson.E{Key: "storage_media", Value: req.StorageMedia})
		}

		if req.VaultType != "" {
			andConditions = append(andConditions, bson.E{Key: "vault_type", Value: req.VaultType})
		}

		if req.TypeIdentity != "" {
			andConditions = append(andConditions, bson.E{Key: "type_identity", Value: req.TypeIdentity})
		}

		if len(req.Tags) > 0 {
			andConditions = append(andConditions, bson.E{Key: "tags", Value: bson.D{{Key: "$in", Value: req.Tags}}})
		}

		if len(andConditions) > 0 {
			orConditions = append(orConditions, andConditions)
		}
	}

	if len(orConditions) > 0 {
		filter["$or"] = orConditions
	}

	return filter
}
