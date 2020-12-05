package main

import (
	"fmt"
	"os"
	"reflect"
)

type a struct {
	X int
	Y float64
	Z string
}

type b struct {
	F int
	G int
	H string
	I float64
}

func main() {
	x := 100
	xRefl := reflect.ValueOf(&x).Elem()
	xType := xRefl.Type()
	fmt.Printf("type is %s\n", xType)
	A := a{100, 200.12, "Struct a"}
	B := b{1, 2, "Struct b", -1.2}
	var r reflect.Value
	arg := os.Args
	if len(arg) == 1 {
		r = reflect.ValueOf(&A).Elem()
	} else {
		r = reflect.ValueOf(&B).Elem()
	}
	iType := r.Type()
	fmt.Printf("type is %s\n", iType)
	fmt.Printf("%d fields if %s are:\n", r.NumField(), iType)

	for i := 0; i < r.NumField(); i++ {
		fmt.Printf("filed name %s", iType.Field(i).Name)
		fmt.Printf(" type %s", iType.Field(i).Type)
		fmt.Printf("  type %s", r.Field(i).Type())
		fmt.Printf(" value %v\n", r.Field(i).Interface())
	}
}
