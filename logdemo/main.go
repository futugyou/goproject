package main

import (
	"log"
	"os"
)

func main() {
	log.Println("1234")
	log.Fatalln("1234")
	log.Panicln("1234")
}

func init() {
	log.SetPrefix("TODO: ")
	f, err := os.OpenFile("./log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
}
