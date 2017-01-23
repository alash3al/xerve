package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

import (
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

import "github.com/gorilla/handlers"

var (
	ROOT    = flag.String("root", ".", "the document root directory")
	ADDR    = flag.String("addr", "127.0.0.1:80", "the listen address")
	INFO    = flag.Bool("info", true, "whether to set the 'Server' header or not")
	GZIP    = flag.Int("gzip", 7, "gzip level, use '0' to disable it, 9 for maximum compression")
	MINIFY  = flag.Bool("minify", true, "mifnity the static files, currently it supports (css, html, js, xml, json, svg+xml)")
	VERSION = "v1.0.0"
)

func init() {
	flag.Usage = func() {
		fmt.Println("xerve, version " + VERSION + " COPYRIGHT 2017 xerve")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	m := minify.New()

	if *MINIFY {
		m.AddFunc("text/css", css.Minify)
		m.AddFunc("text/html", html.Minify)
		m.AddFunc("image/svg+xml", svg.Minify)
		m.AddFuncRegexp(regexp.MustCompile("[/+]javascript$"), js.Minify)
		m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
		m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	}

	log.Fatal(http.ListenAndServe(
		*ADDR,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if *INFO {
				w.Header().Set("Server", "xerve/"+VERSION)
				w.Header().Set("X-Powered-By", "xerve/"+VERSION)
			}
			handlers.CompressHandlerLevel(
				m.Middleware(http.FileServer(http.Dir(*ROOT))),
				*GZIP,
			).ServeHTTP(w, r)
		}),
	))
}
