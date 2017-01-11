package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"text/template"

	syslog "github.com/hashicorp/go-syslog"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

// TODO: override with flags
var admin = &person{Name: "John Weldon", Email: "johnweldon4@gmail.com", Phone: "503-941-0825"}

func main() {
	app := cli.NewApp()
	app.Usage = "Start squidGuard block landing page"
	app.HideVersion = true
	app.Action = Main
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "syslogProtocol",
			EnvVar: "BLOCK_SYSLOG_PROTOCOL",
			Usage:  "syslog network protocol",
			Value:  "udp",
		},
		cli.StringFlag{
			Name:   "syslogHost",
			EnvVar: "BLOCK_SYSLOG_HOST",
			Usage:  "remote syslog host; if empty will default to local",
			Value:  "",
		},
		cli.IntFlag{
			Name:   "syslogPort",
			EnvVar: "BLOCK_SYSLOG_PORT",
			Usage:  "remote syslog port",
			Value:  514,
		},
		cli.IntFlag{
			Name:   "port, p",
			EnvVar: "BLOCK_PORT",
			Usage:  "website port",
			Value:  9000,
		},
	}
	app.Run(os.Args)
}

func Main(c *cli.Context) error {
	out := make(chan blocked)
	go syslogger(c.String("syslogProtocol"), c.String("syslogHost"), c.Int("syslogPort"), out)
	serve(":9000", out)
	return nil
}

const (
	defaultPriority = syslog.LOG_NOTICE
	defaultFacility = "syslog"
	defaultTag      = "web blocked"
)

func syslogger(network, host string, port int, in <-chan blocked) {
	logger, err := syslog.DialLogger(network, fmt.Sprintf("%s:%d", host, port), defaultPriority, defaultFacility, defaultTag)
	if err != nil {
		fmt.Printf("failed to open remote syslog connection %s://%s:%d: %v\n", network, host, port, err)
		logger, err = syslog.NewLogger(defaultPriority, defaultFacility, defaultTag)
		if err != nil {
			fmt.Printf("failed to open local syslog connection: %v\n", err)
			logger = nil
		}
	}
	for {
		select {
		case b := <-in:
			fmt.Printf("Blocked: %s\n", b)
			if logger != nil {
				if _, err := logger.Write([]byte(b.String())); err != nil {
					fmt.Printf("failed to write to syslog: %v\n", err)
				}
			}

		}
	}
}

func serve(port string, out chan<- blocked) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	n := negroni.New(negroni.NewRecovery(), &logger{out: out}, negroni.NewStatic(http.Dir("public")))
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
	category := query.Get("target")
	original := query.Get("requesturl")
	clientip := query.Get("clientip")
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
