package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	viewPath = "/view/"
)

var mux *http.ServeMux

func init() {
	mux = http.NewServeMux()
	mux.HandleFunc("/", loveHandler)
	mux.HandleFunc(viewPath, viewHandler)
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func loveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) == len(viewPath) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Must provide page name")
		return
	}
	title := r.URL.Path[len(viewPath):]
	p, err := loadPage(title)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			w.WriteHeader(404)
			fmt.Fprint(w, "Could not find page")
			return
		} else {
			w.WriteHeader(500)
			log.Fatal(err)
		}
	}

	fBytes, err := os.ReadFile("templates/page-template.txt")
	if err != nil {
		w.WriteHeader(500)
		log.Fatal(err)
	}
	t, err := template.New("page").Parse(string(fBytes))
	if err != nil {
		w.WriteHeader(500)
		log.Fatal(err)
	}
	if err = t.Execute(w, p); err != nil {
		w.WriteHeader(500)
		log.Fatal(err)
	}
}
