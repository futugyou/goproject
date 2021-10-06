package main

type iclothing interface {
	setColor(color string)
	setSize(size int)
	getColor() string
	getSize() int
}

type clothing struct {
	color string
	size  int
}

func (c *clothing) setColor(color string) {
	c.color = color
}
func (c *clothing) setSize(size int) {
	c.size = size
}
func (c *clothing) getColor() string {
	return c.color
}
func (c *clothing) getSize() int {
	return c.size
}
