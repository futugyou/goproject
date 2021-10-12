package main

type icommand interface {
	exec()
}

type oncommand struct {
	device device
}

func (c *oncommand) exec() {
	c.device.on()
}

type offcommand struct {
	device device
}

func (c *offcommand) exec() {
	c.device.off()
}
