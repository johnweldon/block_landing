package main

import (
	"net/http"
	"net/url"
	"text/template"

	"github.com/urfave/negroni"
)

// TODO: override with flags
var admin = &person{Name: "John Weldon", Email: "johnweldon4@gmail.com", Phone: "503-941-0825"}

func main() {
	serve(":9000")
}

func serve(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.Dir("public")))
	n.UseHandler(mux)
	n.Run(port)
}

func index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index.html").ParseFiles("templates/index.html"))
	t.Execute(w, newBlock(r, admin))
}

type person struct {
	Name  string
	Email string
	Phone string
}

type block struct {
	Category    string
	OriginalURL string
	ClientIP    string
	QueryParams map[string][]string
	URL         *url.URL
	Request     *http.Request
	Admin       *person
}

func newBlock(r *http.Request, admin *person) block {
	url := r.URL
	query := r.URL.Query()
	category := ""
	original := ""
	clientip := ""
	if len(query["target"]) > 0 {
		category = query["target"][0]
	}
	if len(query["uri"]) > 0 {
		original = query["uri"][0]
	}
	if len(query["clientip"]) > 0 {
		clientip = query["clientip"][0]
	}
	return block{
		Category:    category,
		OriginalURL: original,
		ClientIP:    clientip,
		QueryParams: query,
		URL:         url,
		Request:     r,
		Admin:       admin,
	}
}
