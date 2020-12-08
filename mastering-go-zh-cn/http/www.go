package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "serving: %s\n", r.URL.Path)
	fmt.Printf("served: %s\n", r.Host)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.RFC1123)
	body := "the current time is:"
	fmt.Fprintf(w, "<h1 align=\"center\">%s</h1>", body)
	fmt.Fprintf(w, "<h2 align=\"center\">%s</h2>\n", t)
	fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Fprintf(w, "Served time for: %s\n", r.Host)
}

type myElement struct {
	Name    string
	Surname string
	Id      string
}

var DATA = make(map[string]myElement)
var DATAFILE = "/tmp/dataFile.gob"

func save() error {
	fmt.Println("saving", DATAFILE)
	err := os.Remove(DATAFILE)
	if err != nil {
		fmt.Println(err)
	}
	saveTo, err := os.Create(DATAFILE)
	if err != nil {
		fmt.Println("cannot create", DATAFILE)
		return err
	}
	defer saveTo.Close()

	encoder := gob.NewEncoder(saveTo)
	err = encoder.Encode(DATA)
	if err != nil {
		fmt.Println("cannot save to", DATAFILE)
		return err
	}
	return nil
}

func load() error {
	fmt.Println("Loading", DATAFILE)
	loadfrom, err := os.Open(DATAFILE)
	loadfrom.Close()
	if err != nil {
		fmt.Println("empty key/value store")
		return err
	}

	decode := gob.NewDecoder(loadfrom)
	decode.Decode(&DATA)
	return nil
}

func Add(k string, n myElement) bool {
	if k == "" {
		return false
	}
	if Lookup(k) == nil {
		DATA[k] = n
		return true
	}
	return false
}

func Delete(k string) bool {
	if Lookup(k) != nil {
		delete(DATA, k)
		return true
	}
	return false
}

func Lookup(k string) *myElement {
	n, ok := DATA[k]
	if ok {
		return &n
	} else {
		return nil
	}
}

func Change(k string, n myElement) bool {
	DATA[k] = n
	return true
}

func Print() {
	for k, v := range DATA {
		fmt.Printf("key: %s, value: %v\n", k, v)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving", r.Host, "for", r.URL.Path)
	myT := template.Must(template.ParseGlob("home.gohtml"))
	myT.ExecuteTemplate(w, "home.gohtml", nil)
}

func listAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing the contents of the KV store!")
	fmt.Fprintf(w, "<a href=\"/\" style=\"margin-right: 20px;\">Home sweet home!</a>")
	fmt.Fprintf(w, "<a href=\"/list\" style=\"margin-right: 20px;\">List all elements!</a>")
	fmt.Fprintf(w, "<a href=\"/change\" style=\"margin-right: 20px;\">Change an elements!</a>")
	fmt.Fprintf(w, "<a href=\"/insert\" style=\"margin-right: 20px;\">Insert an elements!</a>")

	fmt.Fprintf(w, "<h1>The contents of the KV store are:</h1>")
	fmt.Fprintf(w, "<ul>")
	for k, v := range DATA {
		fmt.Fprintf(w, "<li>")
		fmt.Fprintf(w, "<strong>%s</strong> with value: %v\n", k, v)
		fmt.Fprintf(w, "</li>")
	}
	fmt.Fprintf(w, "</ul>")
}

func changeElement(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Change an element of the KV store!")
	tmpl := template.Must(template.ParseFiles("update.gohtml"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	key := r.FormValue("key")
	n := myElement{
		Name:    r.FormValue("name"),
		Surname: r.FormValue("surname"),
		Id:      r.FormValue("id"),
	}

	if !Change(key, n) {
		fmt.Println("Update operation failed!")
	} else {
		err := save()
		if err != nil {
			fmt.Println(err)
			return
		}
		tmpl.Execute(w, struct{ Struct bool }{true})
	}
}

func insertElement(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inserting an element to the KV store!")
	tmpl := template.Must(template.ParseFiles("insert.gohtml"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	key := r.FormValue("key")
	n := myElement{
		Name:    r.FormValue("name"),
		Surname: r.FormValue("surname"),
		Id:      r.FormValue("id"),
	}

	if !Add(key, n) {
		fmt.Println("Add operation failed!")
	} else {
		err := save()
		if err != nil {
			fmt.Println(err)
			return
		}
		tmpl.Execute(w, struct{ Success bool }{true})
	}
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/time", timeHandler)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/change", changeElement)
	r.HandleFunc("/list", listAll)
	r.HandleFunc("/insert", insertElement)
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      r,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
