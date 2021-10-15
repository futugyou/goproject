package main

import "fmt"

func main() {
	m := mementomanager{}
	o := option{age: 10, name: "like", imementomanager: &m}
	o.updateage(11)
	o.updatename("unlike")
	fmt.Println(o.name)
	o = o.getlastchange()
	fmt.Println(o.name)
}
