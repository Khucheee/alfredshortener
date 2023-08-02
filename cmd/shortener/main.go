package main

import (
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.SalamPost)
	mux.HandleFunc("/5k2tWE6P7gjtwJPch8rjZo9JcKzNBSeP1FuAQZGK/", app.Salamget)
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
