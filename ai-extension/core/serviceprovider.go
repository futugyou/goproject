package core

import (
	"fmt"
	"reflect"
	"sync"
)

type IServiceProvider interface {
	Register(service any)
	GetService(key any) (any, error)
}

func GetService[T any](sp IServiceProvider) T {
	if sp == nil {
		panic("service provider is nil")
	}

	t := reflect.TypeOf((*T)(nil)).Elem()

	if service, err := sp.GetService(t); err != nil {
		panic(err.Error())

	} else {
		return service.(T)
	}
}

type ServiceProvider struct {
	services sync.Map
}

func (sp *ServiceProvider) Register(service any) {
	t := reflect.TypeOf(service)
	sp.services.Store(t, service)
}

func (sp *ServiceProvider) GetService(key any) (any, error) {
	if key == nil {
		return nil, fmt.Errorf("key is nil")
	}

	if service, ok := sp.services.Load(key); ok {
		return service, nil
	}

	return nil, fmt.Errorf("service %v not registered", key)
}
