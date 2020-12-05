package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
		savedata()
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

func savedata() {
	s := []byte("data to write\n")
	ff1, err := os.Create("f1.txt")
	if err != nil {
		fmt.Println("cannot create file:", err)
		return
	}
	defer ff1.Close()
	fmt.Fprintf(ff1, string(s))

	ff2, err := os.Create("f2.txt")
	if err != nil {
		fmt.Println("cannot create file ", err)
		return
	}
	defer ff2.Close()
	n, err := ff2.WriteString(string(s))
	fmt.Printf("wrote %d bytes\n", n)

	ff3, err := os.Create("f3.txt")
	if err != nil {
		fmt.Println("cannot create file ", err)
		return
	}
	defer ff3.Close()
	w := bufio.NewWriter(ff3)
	n2, err := w.WriteString(string(s))
	fmt.Printf("wrote %d bytes\n", n2)
	w.Flush()

	ff4 := "f4.txt"
	err = ioutil.WriteFile(ff4, s, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	ff5, err := os.Create("f5.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	n3, err := io.WriteString(ff5, string(s))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("wrote %d bytes\n", n3)
}
