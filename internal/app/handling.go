package app

import (
	"github.com/btcsuite/btcutil/base58"
	"io"
	"net/http"
)

func SalamPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusCreated)
		Body, _ := io.ReadAll(r.Body)
		Bodys := base58.Encode(Body)
		w.Write([]byte("http://localhost:8080/" + Bodys))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
func Salamget(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("location", "https://practicum.yandex.ru/")
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}
