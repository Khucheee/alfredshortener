package main

import (
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	config := app.NewConfig()
	config.SetConfig()
	storage := app.NewStorage()
	logger := app.NewLogger()
	logger.CreateSuggarLogger()
	controller := app.NewBaseController(*config, *storage, *logger)
	r := chi.NewRouter()
	r.Mount("/", controller.Route())

	err := http.ListenAndServe(config.Host, r)
	if err != nil {
		panic(err)
	}
}
