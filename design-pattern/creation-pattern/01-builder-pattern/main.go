package main

import "fmt"

func main() {
	cabinbuilder := getBuilder("cabin")
	igloobuilder := getBuilder("igloo")

	director := newDirector(cabinbuilder)
	cabinhouse := director.builderHouse()

	fmt.Println("cabin house window type: ", cabinhouse.windowType)

	director.setBuilder(igloobuilder)
	igloohouse := director.builderHouse()
	fmt.Println("igloo house window type: ", igloohouse.windowType)

}
