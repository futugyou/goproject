package mongo2struct

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	DBName        string
	ConnectString string
}

type Struct struct {
	Name string
	Item []StructItem
}

type StructItem struct {
	Name string
	Type string
	Tag  string
}

func (m MongoDBConfig) Check() error {
	if len(m.DBName) == 0 {
		return fmt.Errorf("mongodb name can not be nil")
	}
	if len(m.ConnectString) == 0 {
		return fmt.Errorf("mongodb url can not be nil")
	}
	return nil
}

func (m *MongoDBConfig) Generator() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.ConnectString))
	if err != nil {
		log.Println(err)
		return
	}
	db := client.Database(m.DBName)
	filter := bson.D{}
	tables, err := db.ListCollectionSpecifications(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return
	}

	for _, c := range tables {
		generatorTabel(db, c.Name)
	}
}

func generatorTabel(db *mongo.Database, collectionName string) {
	c := db.Collection(collectionName)
	result := c.FindOne(context.Background(), bson.D{})
	b, err := result.Raw()
	if err != nil {
		log.Println(err)
		return
	}
	e, err := b.Elements()
	if err != nil {
		log.Println(err)
		return
	}

	s := Struct{Name: collectionName, Item: make([]StructItem, 0)}

	for _, v := range e {
		itemType := convertBsontypeTogotype(v.Value())
		s.Item = append(s.Item, StructItem{Name: v.Key(), Type: itemType, Tag: fmt.Sprintf("`bson:\"%s\"`", v.Key())})
	}
	fmt.Println(s)
}

func convertBsontypeTogotype(value bson.RawValue) string {
	switch value.Type {
	case bson.TypeDouble:
		return "float64"
	case bson.TypeString:
		return "string"
	case bson.TypeEmbeddedDocument:
	case bson.TypeArray:
		e, err := value.Array().Elements()
		if err != nil || len(e) == 0 {
			return "[]interface{}"
		}
		return "[]" + convertBsontypeTogotype(e[0].Value())
	case bson.TypeBinary:
	case bson.TypeUndefined:
	case bson.TypeObjectID:
		fmt.Println(5)
	case bson.TypeBoolean:
		return "bool"
	case bson.TypeDateTime:
		return "time.Time"
	case bson.TypeNull:
	case bson.TypeRegex:
	case bson.TypeDBPointer:
	case bson.TypeJavaScript:
	case bson.TypeSymbol:
	case bson.TypeCodeWithScope:
	case bson.TypeInt32:
		return "int32"
	case bson.TypeTimestamp:
		return "int64"
	case bson.TypeInt64:
		return "int64"
	case bson.TypeDecimal128:
		return "float64"
	case bson.TypeMinKey:
	case bson.TypeMaxKey:
	}

	return "interface{}"
}
