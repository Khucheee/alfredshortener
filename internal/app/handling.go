package app

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

func SolvePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	reqBody, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	respBody := *setAddress() + base58.Encode(reqBody)
	w.Write([]byte(respBody))
}
func SolveGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", string(base58.Decode(chi.URLParam(r, "shorturl"))))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
