package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/urfave/negroni"
)

type blocked struct {
	ClientIP         net.IP
	ClientDomain     string
	UserID           string
	RequestURL       string
	SourceGroup      string
	DestinationGroup string
}

func (b *blocked) String() string {
	return fmt.Sprintf(
		"category=%q clientip=%q clientdomain=%q userid=%q originalurl=%q group=%q",
		b.DestinationGroup,
		b.ClientIP.String(),
		b.ClientDomain,
		b.UserID,
		b.RequestURL,
		b.SourceGroup,
	)
}

func defaultURLParse(u *url.URL) blocked {
	vals := u.Query()
	return blocked{
		ClientIP:         net.ParseIP(vals.Get("clientip")),
		ClientDomain:     vals.Get("clientdomain"),
		UserID:           vals.Get("userid"),
		RequestURL:       vals.Get("requesturl"),
		SourceGroup:      vals.Get("source"),
		DestinationGroup: vals.Get("target"),
	}
}

type logger struct {
	out chan<- blocked
}

func (l *logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	if r.URL.Path != "/" {
		fmt.Printf("%s %s %q\n", start.Format(time.RFC3339), r.Method, r.URL.Path)
		next(rw, r)
		return
	}
	if l.out != nil {
		l.out <- defaultURLParse(r.URL)
	}

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	fmt.Printf("completed %v %s in %v\n", res.Status(), http.StatusText(res.Status()), time.Since(start))
}
