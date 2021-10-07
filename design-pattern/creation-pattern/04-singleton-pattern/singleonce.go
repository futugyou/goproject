package main

import (
	"fmt"
	"sync"
)

type singleonce struct{}

var once sync.Once
var singleonceInstance *singleonce

func GetSingleonceInstance() *singleonce {
	if singleonceInstance == nil {
		once.Do(func() {
			fmt.Println("new singleonce instance.")
			singleonceInstance = &singleonce{}
		})
	}
	fmt.Println("got singleonce instance.")
	return singleonceInstance
}
