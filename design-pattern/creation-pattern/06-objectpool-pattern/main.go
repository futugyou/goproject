package main

import (
	"log"
	"strconv"
)

func main() {
	conns := make([]iobjectpool, 0)
	for i := 0; i < 3; i++ {
		c := &connectionobject{connstring: strconv.Itoa(i)}
		conns = append(conns, c)
	}
	pool, err := initpool(conns)
	if err != nil {
		log.Fatalf("init error: %s", err)
	}
	con1, err := pool.getobject()
	if err != nil {
		log.Fatalf("getobject error: %s", err)
	}
	con2, err := pool.getobject()
	if err != nil {
		log.Fatalf("getobject error: %s", err)
	}
	pool.receiveobject(con1)
	pool.receiveobject(con2)
	log.Println("ok")
}
