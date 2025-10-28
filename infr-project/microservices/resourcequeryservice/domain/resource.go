package domain

import (
	"time"

	"github.com/futugyou/domaincore/domain"
)

type Resource struct {
	domain.Aggregate `bson:",inline"`
	Name             string    `bson:"name"`
	Type             string    `bson:"type"`
	Data             string    `bson:"data"`
	ImageData        string    `bson:"imageData"`
	Version          int       `bson:"version"`
	IsDelete         bool      `bson:"is_deleted"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
	Tags             []string  `bson:"tags"`
}

func (r Resource) AggregateName() string {
	return "resources_query"
}
