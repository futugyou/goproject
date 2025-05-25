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

func GetService[T any](sp IServiceProvider) (T, error) {
	var zero T

	if sp == nil {
		return zero, fmt.Errorf("service provider is nil")
	}

	t := reflect.TypeOf((*T)(nil)).Elem()
	service, err := sp.GetService(t)
	if err != nil {
		return zero, err
	}

	s, ok := service.(T)
	if !ok {
		return zero, fmt.Errorf("registered service is not of expected type %v", t)
	}

	return s, nil
}

// TODO, find a DI container
// [dig](https://github.com/uber-go/dig) or [wire](https://github.com/google/wire)
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
