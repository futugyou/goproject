package mongo2struct

import (
	"context"
	"github/go-project/tour/internal/common"
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
	Template        *common.Template
	BasePackageName string
	CoreFoler       string
	MongoRepoFolder string
}

func NewManager(db *mongo.Database, entityFolder string, repoFolder string, pkgName string, coreFoler string, mongoRepoFolder string) *Manager {
	return &Manager{
		DB:              db,
		EntityFolder:    entityFolder,
		RepoFolder:      repoFolder,
		Template:        common.NewDefaultTemplate(base_mongorepo_TplString),
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
	return m.Template.GenerateCore(common.CoreConfig{
		PackageName: m.CoreFoler,
		Folder:      m.CoreFoler,
	})
}

func (m *Manager) getEntityStructList() ([]common.EntityStruct, error) {
	tables, err := m.DB.ListCollectionSpecifications(context.Background(), bson.D{})
	if err != nil {
		return []common.EntityStruct{}, err
	}
	return m.createEntityList(tables)
}

func (m *Manager) generatorEntity(entityList []common.EntityStruct) error {
	var wg sync.WaitGroup
	for _, entity := range entityList {
		wg.Add(1)
		go func(entity common.EntityStruct, wg *sync.WaitGroup) {
			defer wg.Done()
			m.Template.GenerateEntity(entity)
		}(entity, &wg)
	}

	wg.Wait()
	return nil
}

func (m *Manager) createEntityList(tables []*mongo.CollectionSpecification) ([]common.EntityStruct, error) {
	entityList := make([]common.EntityStruct, 0)
	ch := make(chan *common.EntityStruct)
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

func (m *Manager) createEntitySingle(name string, wg *sync.WaitGroup, ch chan *common.EntityStruct) {
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

func (m *Manager) generatorRepository(eList []common.EntityStruct) error {
	// base mongo repository implement
	obj := common.BaseRepoImplConfig{
		Folder:   m.MongoRepoFolder,
		FileName: "respository",
		TemplateObj: struct {
			PackageName     string
			BasePackageName string
		}{
			PackageName:     m.MongoRepoFolder,
			BasePackageName: m.BasePackageName,
		},
	}

	err := m.Template.GenerateBaseRepoImpl(obj)
	if err != nil {
		return err
	}

	// other mongo repository implement
	list := make([]common.RepositoryStruct, 0)
	for _, v := range eList {
		list = append(list, common.RepositoryStruct{
			BasePackageName:      m.BasePackageName,
			FileName:             v.FileName,
			RepoName:             v.StructName,
			PackageName:          m.MongoRepoFolder,
			Folder:               m.MongoRepoFolder,
			InterfacePackageName: m.RepoFolder,
			InterfaceFolder:      m.RepoFolder,
		})
	}
	var wg sync.WaitGroup
	for _, entity := range list {
		wg.Add(1)
		go func(entity common.RepositoryStruct, wg *sync.WaitGroup) {
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
