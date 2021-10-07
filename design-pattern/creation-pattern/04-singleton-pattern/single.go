package main

import (
	"fmt"
	"sync"
)

type single struct{}

var lock = &sync.Mutex{}
var singleInstance *single

func GetSingleInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("new single instance.")
			singleInstance = &single{}
		}
	}
	fmt.Println("got single instance.")
	return singleInstance
}
