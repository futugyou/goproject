package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
)

// go run main.go --name=09 go --name=7655
// 2020/10/02 09:19:23 nameFlag: go flag demo 09
// 2020/10/02 09:19:23 name: 7655

// Value can be set
type Value interface {
	String() string
	Set(string) error
}

// Name imp Value
type Name string

func (i *Name) String() string {
	return fmt.Sprint(*i)
}

// Set the flag
func (i *Name) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("name already set")
	}
	*i = Name("go flag demo " + value)
	return nil
}

var name string

func OrderCodeWtihFlag() {
	// ord code
	var nameFlag Name
	flag.Var(&nameFlag, "name", "help info")
	flag.Parse()
	log.Printf("nameFlag: %s", nameFlag)

	goCmd := flag.NewFlagSet("go", flag.ExitOnError)
	goCmd.StringVar(&name, "name", "go project", "help info")
	phpCmd := flag.NewFlagSet("php", flag.ExitOnError)
	phpCmd.StringVar(&name, "n", "php project", "help info")

	args := flag.Args()
	switch args[0] {
	case "go":
		_ = goCmd.Parse(args[1:])
	case "php":
		_ = phpCmd.Parse(args[1:])
	}

	log.Printf("name: %s", name)
}
