package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		// go GetSingleInstance()
		// go GetSingleinitInstance()
		go GetSingleonceInstance()
	}

	fmt.Scanln()
}
