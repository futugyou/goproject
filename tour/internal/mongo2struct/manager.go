package mongo2struct

import (
	"context"
	"log"
	"os/exec"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Manager struct {
	DB           *mongo.Database
	EntityFolder string
	RepoFolder   string
}

func NewManager(db *mongo.Database, entityFolder string, repoFolder string) *Manager {
	return &Manager{
		DB:           db,
		EntityFolder: entityFolder,
		RepoFolder:   repoFolder,
	}
}

func (m *Manager) Generator() {
	m.generatorEntity()
	m.generatorRepository()
	m.formatCode()
}

func (m *Manager) generatorEntity() error {
	tables, err := m.DB.ListCollectionSpecifications(context.Background(), bson.D{})
	if err != nil {
		log.Println(err)
		return err
	}

	entityList, err := m.createEntityList(tables)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, entity := range entityList {
		t := NewTemplate()
		t.GenerateEntity(entity)
	}

	return nil
}

func (m *Manager) createEntityList(tables []*mongo.CollectionSpecification) ([]EntityStruct, error) {
	entityList := make([]EntityStruct, 0)
	for _, c := range tables {
		eles, err := m.createRawElements(c.Name)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		builder := NewEntityStructBuilder(m.EntityFolder, c.Name, eles)
		entity := builder.Build()
		entityList = append(entityList, *entity)
	}
	return entityList, nil
}

func (m *Manager) createRawElements(name string) ([]bson.RawElement, error) {
	c := m.DB.Collection(name)
	result := c.FindOne(context.Background(), bson.D{})
	b, err := result.Raw()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return b.Elements()
}

func (m *Manager) generatorRepository() {
}

func (m *Manager) formatCode() {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Run()
}
