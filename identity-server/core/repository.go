package core

type IInsertRepository[E IEntity] interface {
	Insert(obj E) error
}

type IUpdateRepository[E IEntity, K any] interface {
	Update(obj E, id K) error
}

type IDeleteRepository[E IEntity, K any] interface {
	Delete(obj E, id K) error
}

type IGetAllRepository[E IEntity] interface {
	GetAll() ([]*E, error)
}

type IGetRepository[E IEntity, K any] interface {
	Get(id K) (*E, error)
}
