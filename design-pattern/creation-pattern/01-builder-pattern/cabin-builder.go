package main

type cabinBuilder struct {
	house
}

func newcabinBuilder() *cabinBuilder {
	return &cabinBuilder{}
}

func (b *cabinBuilder) setWindowType() {
	b.windowType = "a window"
}
func (b *cabinBuilder) setDoorType() {
	b.doorType = "a door"
}
func (b *cabinBuilder) setFloor() {
	b.floor = 1
}
func (b *cabinBuilder) getHouse() house {
	return b.house
}
