package main

import (
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	config := new(app.Configure)
	config.SetConfig()
	controller := app.NewBaseController(*config)
	r := chi.NewRouter()
	r.Mount("/", controller.Route())

	err := http.ListenAndServe(config.Host, r)
	if err != nil {
		panic(err)
	}
}
