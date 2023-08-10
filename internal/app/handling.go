package app

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

type BaseController struct {
	config Configure
	urls   map[string]string //мапа содержит сокращенный урл и полный
}

func (b *BaseController) addURL(shorturl, url string) { //добавляем значение в мапу
	b.urls[shorturl] = url
}
func (b *BaseController) searchURL(shorturl string) string { //ищем значение в мапе, если ""то не нашли
	url := b.urls[shorturl]
	return url
}

func NewBaseController(c Configure) *BaseController {
	return &BaseController{
		config: c,
		urls:   make(map[string]string),
	}
}

func (b *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", b.solvePost)
		r.Get("/{shorturl}", b.solveGet)
	})
	return r
}

func (b *BaseController) solvePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	reqBody, _ := io.ReadAll(r.Body)
	reqBodyEncoded := base58.Encode(reqBody)

	defer r.Body.Close()
	respBody := b.config.Address + reqBodyEncoded
	w.Write([]byte(respBody))
	b.addURL(reqBodyEncoded, string(reqBody))
}
func (b *BaseController) solveGet(w http.ResponseWriter, r *http.Request) {
	if b.searchURL(chi.URLParam(r, "shorturl")) == "" { //если ключ в мапе пустой, то 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", b.urls[chi.URLParam(r, "shorturl")]) //если дошли до сюда, то в location суем значение из мапы по ключу
	w.WriteHeader(http.StatusTemporaryRedirect)
}
