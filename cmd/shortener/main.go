package main

import (
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	config := new(app.Configure) //написть функцию для создания
	config.SetConfig()
	storage := app.Storage{Urls: make(map[string]string)}
	controller := app.NewBaseController(*config, storage)
	r := chi.NewRouter()
	r.Mount("/", controller.Route())

	err := http.ListenAndServe(config.Host, r)
	if err != nil {
		panic(err)
	}
}
