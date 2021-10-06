package main

type ibuilder interface {
	setWindowType()
	setDoorType()
	setFloor()
	getHouse() house
}

func getBuilder(builderType string) ibuilder {
	if builderType == "cabin" {
		return &cabinBuilder{}
	}
	if builderType == "igloo" {
		return &iglooBuilder{}
	}
	return nil
}
