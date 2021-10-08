package main

import "fmt"

type folder struct {
	name string
	list []inode
}

func (f *folder) print(prex string) {
	fmt.Println(prex + f.name)
	for _, v := range f.list {
		v.print(prex + prex)
	}
}

func (f *folder) clone() inode {
	folder := &folder{name: f.name + "_clone"}
	var list []inode
	for _, v := range f.list {
		copy := v.clone()
		list = append(list, copy)
	}
	folder.list = list
	return folder
}
