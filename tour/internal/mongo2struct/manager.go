package mongo2struct

import (
	"context"
	"log"
	"os/exec"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Manager struct {
	DB              *mongo.Database
	EntityFolder    string
	RepoFolder      string
	Template        *Template
	BasePackageName string
	CoreFoler       string
	MongoRepoFolder string
}

func NewManager(db *mongo.Database, entityFolder string, repoFolder string, pkgName string, coreFoler string, mongoRepoFolder string) *Manager {
	return &Manager{
		DB:              db,
		EntityFolder:    entityFolder,
		RepoFolder:      repoFolder,
		Template:        NewTemplate(),
		BasePackageName: pkgName,
		CoreFoler:       coreFoler,
		MongoRepoFolder: mongoRepoFolder,
	}
}

func (m *Manager) Generator() {
	if err := m.generatorCore(); err != nil {
		log.Println(err)
		return
	}

	list, err := m.getEntityStructList()
	if err != nil {
		log.Println(err)
		return
	}
	if err := m.generatorEntity(list); err != nil {
		log.Println(err)
		return
	}
	m.generatorRepository(list)
	m.formatCode()
}

func (m *Manager) generatorCore() error {
	return m.Template.GenerateCore(CoreConfig{
		PackageName: m.CoreFoler,
		Folder:      m.CoreFoler,
	})
}

func (m *Manager) getEntityStructList() ([]EntityStruct, error) {
	tables, err := m.DB.ListCollectionSpecifications(context.Background(), bson.D{})
	if err != nil {
		return []EntityStruct{}, err
	}

	return m.createEntityList(tables)
}

func (m *Manager) generatorEntity(entityList []EntityStruct) error {
	var wg sync.WaitGroup
	for _, entity := range entityList {
		wg.Add(1)
		go func(entity EntityStruct, wg *sync.WaitGroup) {
			defer wg.Done()
			m.Template.GenerateEntity(entity)
		}(entity, &wg)
	}

	wg.Wait()
	return nil
}

func (m *Manager) createEntityList(tables []*mongo.CollectionSpecification) ([]EntityStruct, error) {
	entityList := make([]EntityStruct, 0)
	ch := make(chan *EntityStruct)
	var wg sync.WaitGroup
	for _, c := range tables {
		wg.Add(1)
		go m.createEntitySingle(c.Name, &wg, ch)
	}
	go func() {
		for v := range ch {
			if v != nil {
				entityList = append(entityList, *v)
			}
		}
	}()

	wg.Wait()
	close(ch)
	return entityList, nil
}

func (m *Manager) createEntitySingle(name string, wg *sync.WaitGroup, ch chan *EntityStruct) {
	defer wg.Done()

	eles, err := m.createRawElements(name)
	if err != nil {
		log.Println(err)
		ch <- nil
	}

	builder := NewEntityStructBuilder(m.EntityFolder, name, eles)
	ch <- builder.Build()
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

func (m *Manager) generatorRepository(eList []EntityStruct) error {
	obj := struct {
		BasePackageName string
	}{
		BasePackageName: m.BasePackageName,
	}

	err := m.Template.GenerateBaseRepoImpl(obj)
	if err != nil {
		return err
	}
	list := make([]RepositoryStruct, 0)
	for _, v := range eList {
		list = append(list, RepositoryStruct{
			BasePackageName: m.BasePackageName,
			FileName:        v.FileName,
			RepoName:        v.StructName,
		})
	}
	var wg sync.WaitGroup
	for _, entity := range list {
		wg.Add(1)
		go func(entity RepositoryStruct, wg *sync.WaitGroup) {
			defer wg.Done()
			m.Template.GenerateRepository(entity)
		}(entity, &wg)
	}

	wg.Wait()
	return nil
}

func (m *Manager) formatCode() {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Run()
}
