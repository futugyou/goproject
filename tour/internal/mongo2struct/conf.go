package mongo2struct

import (
	"context"
	"fmt"
	"github/go-project/tour/internal/word"
	"log"
	"os/exec"
	"slices"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	EntityFolder  string
	RepoFolder    string
	DBName        string
	ConnectString string
}

type Struct struct {
	EntityFolder string
	RepoFolder   string
	PackageName  string
	FileName     string
	StructName   string
	Items        []StructItem
	Imports      []string
}

type StructItem struct {
	Name string
	Type string
	Tag  string
}

func (m *MongoDBConfig) Check() error {
	if len(m.DBName) == 0 {
		return fmt.Errorf("mongodb name can not be nil")
	}
	if len(m.ConnectString) == 0 {
		return fmt.Errorf("mongodb url can not be nil")
	}
	if len(m.EntityFolder) == 0 {
		m.EntityFolder = "entity"
	}
	if len(m.RepoFolder) == 0 {
		m.RepoFolder = "repository"
	}
	return nil
}

func (m *MongoDBConfig) FormatCode() {
	cmd := exec.Command("go", "fmt", fmt.Sprintf("./%s", m.EntityFolder))
	cmd.Run()
	cmd = exec.Command("go", "fmt", fmt.Sprintf("./%s", m.RepoFolder))
	cmd.Run()
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
		s, _ := m.generatorStruct(db, c.Name)
		t := NewStructTemplate()
		t.Generate(*s)
	}

	m.FormatCode()
}

func (m *MongoDBConfig) generatorStruct(db *mongo.Database, collectionName string) (*Struct, error) {
	c := db.Collection(collectionName)
	result := c.FindOne(context.Background(), bson.D{})
	b, err := result.Raw()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	e, err := b.Elements()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s := &Struct{
		FileName:     collectionName,
		EntityFolder: m.EntityFolder,
		RepoFolder:   m.RepoFolder,
		PackageName:  m.EntityFolder,
		StructName:   UnderscoreToUpperCamelCase(collectionName),
		Items:        make([]StructItem, 0),
		Imports:      make([]string, 0),
	}

	for _, v := range e {
		itemType := convertBsontypeTogotype(v.Value())
		s.Imports = createImports(s.Imports, itemType)
		s.Items = append(s.Items, StructItem{
			Name: UnderscoreToUpperCamelCase(v.Key()),
			Type: itemType,
			Tag:  fmt.Sprintf("`bson:\"%s\"`", v.Key()),
		})
	}

	return s, nil
}

func UnderscoreToUpperCamelCase(s string) string {
	s = word.UnderscoreToUpperCamelCase(s)
	if unicode.IsDigit(rune(s[0])) {
		s = "A" + s
	}
	return s
}

func createImports(s []string, itemType string) []string {
	// for now 'time' only
	if itemType == "time.Time" && !slices.Contains(s, "\"time\"") {
		s = append(s, "\"time\"")
	}

	return s
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
		return "string"
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
