package main

import (
	"fmt"
)

type singleinit struct{}

var singleinitInstance *singleinit

func init() {
	fmt.Println("new singleinit instance.")
	singleinitInstance = &singleinit{}
}

func GetSingleinitInstance() *singleinit {
	if singleinitInstance == nil {
		panic("singleinit is null")
	}
	fmt.Println("got singleinit instance.")
	return singleinitInstance
}
