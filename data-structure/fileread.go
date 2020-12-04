package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

func lineByLine(file string) error {
	var err error
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("read error %s", err)
			break
		}
		fmt.Printf(line)
	}
	return nil
}

func wordByWord(file string) error {
	var err error
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error %s", err)
			return err
		}
		r := regexp.MustCompile("[^\\s]+")
		words := r.FindAllString(line, -1)
		for i := 0; i < len(words); i++ {
			fmt.Printf(words[i])
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("no file")
		return
	}

	for _, file := range flag.Args() {
		err := lineByLine(file)
		if err != nil {
			fmt.Println(err)
		}
	}
	for _, file := range flag.Args() {
		err := wordByWord(file)
		if err != nil {
			fmt.Println(err)
		}
	}
}
