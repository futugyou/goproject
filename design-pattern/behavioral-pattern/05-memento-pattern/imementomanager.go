package main

type imementomanager interface {
	add(m imemento)
	getlast() imemento
}

type mementomanager struct {
	imementolist []imemento
}

func (manager *mementomanager) add(m imemento) {
	manager.imementolist = append(manager.imementolist, m)
}

func (manager *mementomanager) getlast() imemento {
	lenght := len(manager.imementolist)
	if lenght > 1 {
		return manager.imementolist[lenght-2]
	}
	return nil
}
