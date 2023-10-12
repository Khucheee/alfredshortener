package main

import (
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	config := app.NewConfig()
	config.SetConfig()
	//в нью кипер кидаю конфиг и внутри выбираю какой кипер возвращать
	//ctx := context.Background()

	keeper := app.NewKeeper(*config)
	storage := app.NewStorage(*keeper)
	//work := worker.NewWorker(storage)
	logger := app.NewLogger()
	logger.CreateSuggarLogger()
	controller := app.NewBaseController(*config, *storage, *logger)
	//work.start(ctx)
	r := chi.NewRouter()
	r.Mount("/", controller.Route())

	err := http.ListenAndServe(config.Host, r)
	if err != nil {
		panic(err)
	}
}
