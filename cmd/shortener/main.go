package main

import (
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"net/http"
)

const slash string = "/"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.SalamPost)
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
