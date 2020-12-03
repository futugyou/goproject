package main

import (
	"container/ring"
	"fmt"
)

var size int = 10

func main() {
	myring := ring.New(size + 1)
	fmt.Println("empty : ", *myring)
	for i := 0; i < myring.Len()-1; i++ {
		myring.Value = i
		myring = myring.Next()
	}
	myring.Value = 2
	sum := 0
	myring.Do(func(x interface{}) {
		t := x.(int)
		sum += t
	})
	fmt.Println("sum : ", sum)
	for i := 0; i < myring.Len()+12; i++ {
		myring = myring.Next()
		fmt.Print(myring.Value, " ")
	}
}
