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
	editPath = "/edit/"
	savePath = "/save/"
)

var mux *http.ServeMux

func init() {
	mux = http.NewServeMux()
	mux.HandleFunc("/", loveHandler)
	mux.HandleFunc(viewPath, viewHandler)
	mux.HandleFunc(editPath, editHandler)
	mux.HandleFunc(savePath, saveHandler)
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
			http.Redirect(w, r, editPath+title, http.StatusFound)
		} else {
			w.WriteHeader(500)
			log.Print(err)
			return
		}
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(editPath):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		w.WriteHeader(500)
		log.Fatal(err)
	}
	t.Execute(w, p)
}
