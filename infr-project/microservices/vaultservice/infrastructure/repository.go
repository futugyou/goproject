package infrastructure

import (
	"context"
	"fmt"

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
	values := make([]any, len(ids))
	for i, v := range ids {
		values[i] = v
	}

	query := domaincore.NewQuery().
		In("id", values...).
		Build()
	condition := domaincore.NewQueryOptions(nil, nil, nil, query)
	return r.Find(ctx, condition)
}

func (r *VaultRepository) SearchVaults(ctx context.Context, query domain.VaultQuery) ([]domain.Vault, error) {
	filter := buildSearchFilter(query.Filters)
	condition := domaincore.NewQueryOptions(query.Page, query.Size, nil, filter)
	return r.Find(ctx, condition)
}

func buildSearchFilter(reqs []domain.VaultFilter) domaincore.FilterExpr {
	if len(reqs) == 0 {
		return nil
	}

	var orConditions []domaincore.FilterExpr

	for _, req := range reqs {
		var andConditions []domaincore.FilterExpr

		if req.ID != "" {
			andConditions = append(andConditions, domaincore.Eq{
				Field: "id",
				Value: req.ID,
			})
		}

		if req.Description != "" {
			andConditions = append(andConditions, domaincore.Like{
				Field:           "description",
				Pattern:         req.Description,
				CaseInsensitive: true,
			})
		}

		if req.Key != "" {
			if req.KeyFuzzy {
				andConditions = append(andConditions, domaincore.Like{
					Field:           "key",
					Pattern:         req.Key,
					CaseInsensitive: true,
				})
			} else {
				andConditions = append(andConditions, domaincore.Eq{
					Field: "key",
					Value: req.Key,
				})
			}
		}

		if req.StorageMedia != "" {
			andConditions = append(andConditions, domaincore.Eq{
				Field: "storage_media",
				Value: req.StorageMedia,
			})
		}

		if req.VaultType != "" {
			andConditions = append(andConditions, domaincore.Eq{
				Field: "vault_type",
				Value: req.VaultType,
			})
		}

		if req.TypeIdentity != "" {
			andConditions = append(andConditions, domaincore.Eq{
				Field: "type_identity",
				Value: req.TypeIdentity,
			})
		}

		if len(req.Tags) > 0 {
			values := make([]any, len(req.Tags))
			for i, v := range req.Tags {
				values[i] = v
			}
			andConditions = append(andConditions, domaincore.In{
				Field:  "tags",
				Values: values,
			})
		}

		if len(andConditions) == 1 {
			orConditions = append(orConditions, andConditions[0])
		} else if len(andConditions) > 1 {
			orConditions = append(orConditions, domaincore.And(andConditions))
		}
	}

	if len(orConditions) == 0 {
		return nil
	} else if len(orConditions) == 1 {
		return orConditions[0]
	} else {
		return domaincore.Or(orConditions)
	}
}
