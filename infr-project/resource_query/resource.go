package resourcequery

import (
	"time"

	"github.com/futugyou/infr-project/domain"
)

type Resource struct {
	domain.Aggregate `json:",inline" bson:",inline"`
	Name             string    `json:"name" bson:"name"`
	Type             string    `json:"type" bson:"type"`
	Data             string    `json:"data" bson:"data"`
	Version          int       `json:"version" bson:"version"`
	IsDelete         bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
	Tags             []string  `json:"tags" bson:"tags"`
}

func (r Resource) AggregateName() string {
	return "resources_query"
}
