package main

import (
	"net/http"
	"net/url"
	"text/template"

	"github.com/codegangsta/negroni"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":9000")
}

func index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index").Parse(indexTemplate))
	t.Execute(w, NewBlock(r, DefaultAdmin))
}

type Person struct {
	Name  string
	Email string
	Phone string
}

type Block struct {
	Category    string
	OriginalURL string
	ClientIP    string
	QueryParams map[string][]string
	URL         *url.URL
	Request     *http.Request
	Admin       *Person
}

func NewBlock(r *http.Request, admin *Person) Block {
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
	return Block{
		Category:    category,
		OriginalURL: original,
		ClientIP:    clientip,
		QueryParams: query,
		URL:         url,
		Request:     r,
		Admin:       admin,
	}
}

var DefaultAdmin = &Person{Name: "John Weldon", Email: "johnweldon4@gmail.com", Phone: "503-941-0825"}

const indexTemplate = `
<html>

<head>
<title>Intercept {{ .OriginalURL }}</title>
<link href='app.css' rel='stylesheet' />
<script src='app.js'></script>
</head>

<body>

<div id='main'>
<div class='header'></div>
<h1>Blocked</h1>
<p>The url you requested, <span class='url'>{{ .OriginalURL }}</span>, has been deemed inappropriate for this network,
because it is classified as <span class='category'>{{ .Category }}</span>.</p>

{{ with .Admin }}
<p>If you would like to request a review, or an exception, please contact:<br/>
<a href='mailto:{{ .Email }}'>{{ .Name }}</a><br/>
<a href='tel:{{ .Phone }}'>{{ .Phone }}</a>
</p>
{{ end }}

</div>

<div id='raw'>
{{ .Request.URL }}

{{ printf "%#v" .Request }}
</div>

</body>
</html>
`
