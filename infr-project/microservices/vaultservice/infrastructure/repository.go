package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/vaultservice/domain"

	domaincore "github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/mongoimpl"
)

type VaultRepository struct {
	mongoimpl.BaseRepository[domain.Vault]
}

func NewVaultRepository(client *mongo.Client, config mongoimpl.DBConfig) *VaultRepository {
	return &VaultRepository{
		BaseRepository: *mongoimpl.NewBaseRepository[domain.Vault](client, config),
	}
}

func (r *VaultRepository) InsertMultipleVault(ctx context.Context, vaults []domain.Vault) error {
	if len(vaults) == 0 {
		return fmt.Errorf("not data need insert")
	}

	c := r.Client.Database(r.DBName).Collection(vaults[0].AggregateName())
	documents := make([]any, 0)
	for i := 0; i < len(vaults); i++ {
		documents = append(documents, vaults[i])
	}
	_, err := c.InsertMany(ctx, documents)
	return err
}

func (r *VaultRepository) GetVaultByIds(ctx context.Context, ids []string) ([]domain.Vault, error) {
	condition := domaincore.NewQueryOptions(nil, nil, nil, map[string]any{"id": bson.M{"$in": ids}})
	return r.BaseRepository.Find(ctx, condition)
}

func (r *VaultRepository) SearchVaults(ctx context.Context, req []domain.VaultSearch, page *int, size *int) ([]domain.Vault, error) {
	filter := buildSearchFilter(req)
	condition := domaincore.NewQueryOptions(page, size, nil, filter)
	return r.BaseRepository.Find(ctx, condition)
}

func buildSearchFilter(reqs []domain.VaultSearch) map[string]any {
	filter := map[string]any{}
	if len(reqs) == 0 {
		return filter
	}

	orConditions := bson.A{}
	for _, req := range reqs {
		andConditions := bson.D{}

		if req.ID != "" {
			andConditions = append(andConditions, bson.E{Key: "id", Value: req.ID})
		}

		if req.Description != "" {
			andConditions = append(andConditions, bson.E{Key: "description", Value: bson.D{{Key: "$regex", Value: req.Description}, {Key: "$options", Value: "i"}}})
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
