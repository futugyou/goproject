package main

func main() {
	c := &context{}

	var s istrategy = &strategyA{}
	c.setstrategy(s)
	c.exec()

	s = &strategyB{}
	c.setstrategy(s)
	c.exec()

	s = &strategyC{}
	c.setstrategy(s)
	c.exec()
}
