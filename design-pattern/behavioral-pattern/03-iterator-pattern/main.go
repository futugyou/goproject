package main

import "fmt"

func main() {
	u1 := &user{name: "a", age: 1}
	u2 := &user{name: "b", age: 2}

	users := &userCollection{
		users: []*user{u1, u2},
	}

	iterator := users.createiterator()

	for iterator.hasnext() {
		user := iterator.getnext()
		fmt.Printf("User is %+v\n", user)
	}
}
