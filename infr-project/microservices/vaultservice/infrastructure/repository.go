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
	mongoimpl.BaseCRUD[domain.Vault]
}

func NewVaultRepository(client *mongo.Client, config mongoimpl.DBConfig) *VaultRepository {
	if config.CollectionName == "" {
		config.CollectionName = "vaults"
	}

	getID := func(a domain.Vault) string { return a.AggregateId() }

	return &VaultRepository{
		BaseCRUD: *mongoimpl.NewBaseCRUD(client, config, getID),
	}
}

func (r *VaultRepository) InsertMultipleVault(ctx context.Context, vaults []domain.Vault) error {
	if len(vaults) == 0 {
		return fmt.Errorf("not data need insert")
	}

	c := r.Client.Database(r.DBName).Collection(r.CollectionName)
	documents := make([]any, 0)
	for i := range vaults {
		documents = append(documents, vaults[i])
	}
	_, err := c.InsertMany(ctx, documents)
	return err
}

func (r *VaultRepository) GetVaultByIds(ctx context.Context, ids []string) ([]domain.Vault, error) {
	condition := domaincore.NewQueryOptions(nil, nil, nil, map[string]any{"id": bson.M{"$in": ids}})
	return r.Find(ctx, condition)
}

func (r *VaultRepository) SearchVaults(ctx context.Context, query domain.VaultQuery) ([]domain.Vault, error) {
	filter := buildSearchFilter(query.Filters)
	condition := domaincore.NewQueryOptions(query.Page, query.Size, nil, filter)
	return r.Find(ctx, condition)
}

func buildSearchFilter(reqs []domain.VaultFilter) map[string]any {
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
