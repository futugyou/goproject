package main

type context struct {
	strategy istrategy
}

func (c *context) setstrategy(strategy istrategy) {
	c.strategy = strategy
}
func (c *context) exec() {
	c.strategy.exec()
}
