package mongo2struct

const repoInterfaceTplString = `
package {{ .InterfacePackageName }}

import (
	"{{ .BasePackageName }}/core"
	"{{ .BasePackageName }}/entity"
)

type I{{ .RepoName }}Repository interface {
	core.IRepository[entity.{{ .RepoName }}Entity, string]
}
`
const repoMongoImplTplString = `
package {{ .PackageName }}

import (
	"{{ .BasePackageName }}/entity"
)

type {{ .RepoName }}Repository struct {
	*MongoRepository[entity.{{ .RepoName }}Entity, string]
}

func New{{ .RepoName }}Repository(config DBConfig) *{{ .RepoName }}Repository {
	baseRepo := NewMongoRepository[entity.{{ .RepoName }}Entity, string](config)
	return &{{ .RepoName }}Repository{baseRepo}
}

`
