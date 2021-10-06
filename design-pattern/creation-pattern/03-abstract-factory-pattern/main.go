package main

import "fmt"

func main() {
	adsfactory, _ := newigarmentfactory("ads")
	adsclothing := adsfactory.createClothing()
	adspants := adsfactory.createPants()

	fmt.Println(adsclothing.getColor())
	fmt.Println(adsclothing.getSize())
	fmt.Println(adspants.getMaterial())

	nikefactory, _ := newigarmentfactory("nike")
	nikeclothing := nikefactory.createClothing()
	nikepants := nikefactory.createPants()

	fmt.Println(nikeclothing.getColor())
	fmt.Println(nikeclothing.getSize())
	fmt.Println(nikepants.getMaterial())
}
