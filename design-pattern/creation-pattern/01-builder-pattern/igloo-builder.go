package main

type iglooBuilder struct {
	house
}

func newiglooBuilder() *iglooBuilder {
	return &iglooBuilder{}
}

func (b *iglooBuilder) setWindowType() {
	b.windowType = "b window"
}
func (b *iglooBuilder) setDoorType() {
	b.doorType = "b door"
}
func (b *iglooBuilder) setFloor() {
	b.floor = 11
}
func (b *iglooBuilder) getHouse() house {
	return b.house
}
