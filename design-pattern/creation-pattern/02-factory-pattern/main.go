package main

import "fmt"

func main() {
	ak47, _ := getGun("ak47")
	m16, _ := getGun("m16")

	fmt.Println(ak47.getName())
	fmt.Println(ak47.getPower())
	fmt.Println(m16.getName())
	fmt.Println(m16.getPower())
}
