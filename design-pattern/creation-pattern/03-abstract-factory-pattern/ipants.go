package main

type ipants interface {
	setMaterial(material string)
	getMaterial() string
}

type pants struct {
	material string
}

func (p *pants) setMaterial(material string) {
	p.material = material
}

func (p *pants) getMaterial() string {
	return p.material
}
