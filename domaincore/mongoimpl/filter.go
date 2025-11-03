package mongoimpl

import (
	"github.com/futugyou/domaincore/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func buildMongoFilter(expr domain.FilterExpr) bson.D {
	switch e := expr.(type) {
	case domain.Eq:
		return bson.D{{Key: e.Field, Value: e.Value}}
	case domain.Ne:
		return bson.D{{Key: e.Field, Value: bson.D{{Key: "$ne", Value: e.Value}}}}
	case domain.Gt:
		return bson.D{{Key: e.Field, Value: bson.D{{Key: "$gt", Value: e.Value}}}}
	case domain.Gte:
		return bson.D{{Key: e.Field, Value: bson.D{{Key: "$gte", Value: e.Value}}}}
	case domain.Lt:
		return bson.D{{Key: e.Field, Value: bson.D{{Key: "$lt", Value: e.Value}}}}
	case domain.Lte:
		return bson.D{{Key: e.Field, Value: bson.D{{Key: "$lte", Value: e.Value}}}}
	case domain.In:
		return bson.D{{Key: e.Field, Value: bson.D{{Key: "$in", Value: e.Values}}}}
	case domain.Nin:
		return bson.D{{Key: e.Field, Value: bson.D{{Key: "$nin", Value: e.Values}}}}
	case domain.Like:
		regex := bson.D{{Key: "$regex", Value: e.Pattern}}
		if e.CaseInsensitive {
			regex = append(regex, bson.E{Key: "$options", Value: "i"})
		}
		return bson.D{{Key: e.Field, Value: regex}}
	case domain.And:
		arr := bson.A{}
		for _, sub := range e {
			arr = append(arr, buildMongoFilter(sub))
		}
		return bson.D{{Key: "$and", Value: arr}}
	case domain.Or:
		arr := bson.A{}
		for _, sub := range e {
			arr = append(arr, buildMongoFilter(sub))
		}
		return bson.D{{Key: "$or", Value: arr}}
	case domain.Not:
		return bson.D{{Key: "$not", Value: buildMongoFilter(e.Expr)}}
	default:
		return bson.D{}
	}
}
