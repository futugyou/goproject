package common

type EntityStruct struct {
	EntityFolder string
	FileName     string
	PackageName  string
	Imports      []string
	StructName   string
	Items        []EntityStructItem
}

type EntityStructItem struct {
	Name string
	Type string
	Tag  string
}

type RepositoryStruct struct {
	BasePackageName      string
	FileName             string
	RepoName             string
	PackageName          string
	Folder               string
	InterfacePackageName string
	InterfaceFolder      string
}

type CoreConfig struct {
	PackageName string
	Folder      string
}

type BaseRepoImplConfig struct {
	Folder      string
	FileName    string
	TemplateObj interface{}
}
