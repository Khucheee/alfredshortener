package main

import (
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	config := app.NewConfig()
	config.SetConfig()
	keeper := app.NewKeeper(config.FilePath)
	storage := app.NewStorage(keeper)
	for {
		if a, _ := keeper.Restore(); a == "" {
			break
		}
		storage.AddURL(keeper.Restore())
	}
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
