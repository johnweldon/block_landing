package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

type logger struct {
}

func (l *logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = "0.0.0.0"
	}
	ipv := net.ParseIP(ip)
	fwd := r.Header.Get("x-forwarded-for")
	fmt.Printf("\n%s IP: %s %s\n", time.Now().Format(time.RFC3339), ipv, fwd)
	fmt.Printf("  started %s %s\n", r.Method, r.URL.String())

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	fmt.Printf("completed %v %s in %v\n", res.Status(), http.StatusText(res.Status()), time.Since(start))
}
