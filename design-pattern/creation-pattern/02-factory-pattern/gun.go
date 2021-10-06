package main

type gun struct {
	name  string
	power int
}

func (g *gun) getName() string {
	return g.name
}

func (g *gun) setName(name string) {
	g.name = name
}

func (g *gun) getPower() int {
	return g.power
}

func (g *gun) setPower(power int) {
	g.power = power
}
