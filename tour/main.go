package main

import (
	"flag"
	"log"
)

var name string

func main() {
	flag.Parse()
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
