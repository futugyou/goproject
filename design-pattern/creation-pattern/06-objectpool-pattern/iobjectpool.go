package main

type iobjectpool interface {
	getobject() interface{}
}
