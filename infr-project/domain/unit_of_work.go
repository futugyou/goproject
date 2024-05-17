package domain

type IUnitOfWork interface {
	Commit() error
	Rollback() error
}
