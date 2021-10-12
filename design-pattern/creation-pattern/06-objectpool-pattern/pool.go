package main

import (
	"fmt"
	"sync"
)

type pool struct {
	idle   []iobjectpool
	active []iobjectpool
	cap    int
	lock   *sync.Mutex
}

func initpool(poolobjects []iobjectpool) (*pool, error) {
	if len(poolobjects) == 0 {
		return nil, fmt.Errorf("lenght error")
	}
	active := make([]iobjectpool, 0)
	pool := &pool{
		idle:   poolobjects,
		active: active,
		cap:    len(poolobjects),
		lock:   new(sync.Mutex),
	}
	return pool, nil
}

func (p *pool) getobject() (iobjectpool, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.idle) == 0 {
		return nil, fmt.Errorf("no object")
	}
	obj := p.idle[0]
	p.idle = p.idle[1:]
	p.active = append(p.active, obj)
	return obj, nil
}

func (p *pool) removeobject(obj iobjectpool) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	currlen := len(p.active)
	for i, o := range p.active {
		if o.getobject() == obj.getobject() {
			p.active[currlen-1], p.active[i] = p.active[i], p.active[currlen-1]
			p.active = p.active[:currlen]
			return nil
		}
	}
	return fmt.Errorf("not exist")
}

func (p *pool) receiveobject(obj iobjectpool) error {
	err := p.removeobject(obj)
	if err != nil {
		return err
	}
	p.idle = append(p.idle, obj)
	return nil
}
