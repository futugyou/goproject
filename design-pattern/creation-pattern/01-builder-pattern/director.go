package main

type director struct {
	builder ibuilder
}

func newDirector(b ibuilder) *director {
	return &director{
		builder: b,
	}
}

func (d *director) setBuilder(b ibuilder) {
	d.builder = b
}

func (d *director) builderHouse() house {
	d.builder.setWindowType()
	d.builder.setDoorType()
	d.builder.setFloor()
	return d.builder.getHouse()
}
