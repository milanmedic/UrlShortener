package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"urlshortener.mmedic.com/v2/urlshort"
)

func main() {

	filename := flag.String("file", "urls.json", "File that contains a list of url mappings.")
	flag.Parse()

	mux := defaultMux()

	var urlMappings map[string]string = make(map[string]string)
	urlMappings["/urlshort-godoc"] = "https://godoc.org/github.com/gophercises/urlshort"

	mapHandler := urlshort.MapHandler(urlMappings, mux)

	yamlHandler := urlshort.YAMLHandler(ReadFile(*filename), mapHandler)

	jsonHandler := urlshort.JSONHandler(ReadFile(*filename), yamlHandler)

	http.ListenAndServe("localhost:3000", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Hello)
	return mux
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func ReadFile(filename string) []byte {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("file.Get err   #%v ", err)
	}

	return contents
}
