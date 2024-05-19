package infrastructure

type MemoryUnitOfWork struct {
}

func NewMemoryUnitOfWork() *MemoryUnitOfWork {
	return &MemoryUnitOfWork{}
}

func (u *MemoryUnitOfWork) Start() error {
	return nil
}

func (u *MemoryUnitOfWork) Commit() error {
	return nil
}

func (u *MemoryUnitOfWork) Rollback() error {
	return nil
}
