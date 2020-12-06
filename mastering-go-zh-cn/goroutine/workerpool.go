package main

import (
	"fmt"
	"sync"
	"time"
)

type Client struct {
	id      int
	integer int
}
type Data struct {
	job    Client
	square int
}

var (
	size    = 10
	clients = make(chan Client, size)
	data    = make(chan Data, size)
)

func worker(w *sync.WaitGroup) {
	for c := range clients {
		squre := c.integer * c.integer
		output := Data{c, squre}
		data <- output
		time.Sleep(time.Second)
	}
	w.Done()
}

func makewp(n int) {
	var w sync.WaitGroup
	for i := 0; i < n; i++ {
		w.Add(1)
		go worker(&w)
	}
	w.Wait()
	close(data)
}

func create(n int) {
	for i := 0; i < n; i++ {
		c := Client{i, i}
		clients <- c
	}
	close(clients)
}

func main() {
	go create(9)
	finished := make(chan interface{})
	go func() {
		for d := range data {
			fmt.Printf("client id: %d\tint: ", d.job.id)
			fmt.Printf("%d\tsquare: %d\n", d.job.integer, d.square)
		}
		finished <- true
	}()
	makewp(6)
	fmt.Printf(": %v\n", <-finished)
}
