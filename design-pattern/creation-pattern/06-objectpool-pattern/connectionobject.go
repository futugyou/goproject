package main

type connectionobject struct {
	connstring string
}

func (c *connectionobject) getobject() interface{} {
	return c.connstring
}
