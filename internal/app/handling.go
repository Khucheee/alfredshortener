package app

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

type BaseController struct {
	config Configure
}

func NewBaseController(c Configure) *BaseController {
	return &BaseController{
		config: c,
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
	defer r.Body.Close()
	respBody := b.config.Address + base58.Encode(reqBody)
	w.Write([]byte(respBody))
}
func (b *BaseController) solveGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", string(base58.Decode(chi.URLParam(r, "shorturl"))))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
