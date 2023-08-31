package main

import (
	"bufio"
	"encoding/json"
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"github.com/go-chi/chi"
	"net/http"
	"os"
)

func main() {
	config := new(app.Configure)
	logger := new(app.Logger)
	config.SetConfig()
	logger.CreateLogger()
	controller := app.NewBaseController(*config, *logger)

	file, err := os.OpenFile(config.FilePath, os.O_RDONLY|os.O_CREATE, 0666)

	text := bufio.NewReader(file)
	var data []byte
	for {
		data, err = text.ReadBytes('\n')
		if len(data) == 0 {
			break
		}
		jon := app.JSONfile{}
		json.Unmarshal(data, &jon)
		controller.Urls[jon.Short_url] = jon.Original_url
	}
	file.Close()
	r := chi.NewRouter()

	r.Mount("/", controller.Route())

	err = http.ListenAndServe(config.Host, r)
	if err != nil {
		panic(err)
	}
}
