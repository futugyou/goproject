package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var CLOSEA = false
var DATA = make(map[int]bool)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func first(min, max int, out chan<- int) {
	for {
		if CLOSEA {
			close(out)
			return
		}
		out <- random(min, max)
	}
}

func second(out chan<- int, in <-chan int) {
	for x := range in {
		fmt.Print(x, " ")
		_, ok := DATA[x]
		if ok {
			CLOSEA = true
		} else {
			DATA[x] = true
			out <- x
		}
	}
	fmt.Println()
	close(out)
}

func third(in <-chan int) {
	var sum int
	sum = 0
	for x2 := range in {
		sum += x2
	}
	fmt.Printf("the sum is %d\n", sum)
}

func gen(min, max int, createNumber chan int, end chan bool) {
	t := time.NewTimer(4 * time.Second)
	for {
		select {
		case createNumber <- rand.Intn(max-min) + min:
		case <-end:
			close(end)
			return
		case <-t.C: //<-time.After(4 * time.Second):
			fmt.Println("\ntime.After()!")
		}
	}
}

func f1(cc chan chan int, f chan bool) {
	c := make(chan int)
	cc <- c
	defer close(c)
	sum := 0
	select {
	case x := <-c:
		for i := 0; i <= x; i++ {
			sum += i
		}
		c <- sum
	case <-f:
		return
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("need two int param")
		os.Exit(1)
	}
	n1, _ := strconv.Atoi(os.Args[1])
	n2, _ := strconv.Atoi(os.Args[2])
	if n1 > n2 {
		fmt.Printf("%d should be smaller than %d\n", n1, n2)
		return
	}

	rand.Seed(time.Now().UnixNano())
	A := make(chan int)
	B := make(chan int)

	go first(n1, n2, A)
	go second(B, A)
	third(B)

	createNumber := make(chan int)
	end := make(chan bool)

	go gen(0, 2*n1, createNumber, end)
	for i := 0; i < n1; i++ {
		fmt.Printf("%d ", <-createNumber)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("exting...")
	end <- true

	cc := make(chan chan int)
	for i := 1; i < n1+1; i++ {
		f := make(chan bool)
		go f1(cc, f)
		ch := <-cc
		ch <- i
		for sum := range ch {
			fmt.Print("SUM(", i, ")=", sum)
		}
		fmt.Println()
		//time.Sleep(time.Second)
		close(f)
	}
}
