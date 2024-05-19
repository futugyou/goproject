package domain

type IUnitOfWork interface {
	Start() error
	Commit() error
	Rollback() error
}
