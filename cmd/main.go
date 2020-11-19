package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/osery/coffee-maker/internal/brew"
	"github.com/osery/coffee-maker/internal/html"
	"github.com/osery/coffee-maker/internal/rest"
	"github.com/osery/coffee-maker/internal/store"
)

var (
	bind  = flag.String("bind", "0.0.0.0:3334", "Bind address to listen on.")
	debug = flag.Bool("debug", false, "Debug mode with human readable logging.")
)

func main() {
	flag.Parse()
	initLogger()

	s := store.NewInMemory()
	b := brew.NewBrewer(s)
	rh := rest.NewRESTHandler(s, b)
	hh := html.NewHTMLHandler(s)
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/coffees", rh.ListCoffees).Methods("GET")
	r.HandleFunc("/coffees", rh.PostCoffee).Methods("POST")
	r.HandleFunc("/coffees/{name}", rh.GetCoffee).Methods("GET")
	r.HandleFunc("/", hh.IndexPage)

	zap.L().With(zap.String("bind", *bind)).Info("Listening.")
	log.Fatal(http.ListenAndServe(*bind, r))
}

func initLogger() {
	var l *zap.Logger
	if *debug {
		l, _ = zap.NewDevelopment()
	} else {
		l, _ = zap.NewProduction()
	}
	zap.ReplaceGlobals(l)
}
