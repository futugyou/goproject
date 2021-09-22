package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func main() {
	log.Println("1234")
	Trace.Println("1234")
	Info.Println("1234")
	Warning.Println("1234")
	Error.Println("1234")
}

func init() {
	file, err := os.OpenFile("./errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	Trace = log.New(ioutil.Discard, "TRACE: ", log.Ldate|log.Ltime|log.Llongfile)
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Llongfile)
	Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
}
