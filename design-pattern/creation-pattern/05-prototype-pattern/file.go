package main

import "fmt"

type file struct {
	name string
}

func (f *file) print(prex string) {
	fmt.Println(prex + f.name)
}

func (f *file) clone() inode {
	return &file{name: f.name + "_clone"}
}
