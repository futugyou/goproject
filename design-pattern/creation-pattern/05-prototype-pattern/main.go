package main

import "fmt"

func main() {
	f1 := &file{name: "f1"}
	f2 := &file{name: "f2"}
	f3 := &file{name: "f3"}
	folder1 := &folder{
		name: "folder1",
		list: []inode{f1, f2, f3},
	}
	folder1.print(" ")
	fmt.Println("----")
	clonefolder := folder1.clone()
	clonefolder.print(" ")
}
