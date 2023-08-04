package app

import (
	"github.com/btcsuite/btcutil/base58"
	"io"
	"net/http"
)

func SolveRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		respBody := "http://localhost:8080/" + base58.Encode(reqBody)
		w.Write([]byte(respBody))
		return
	}
	if r.Method == http.MethodGet {
		shorturl := r.URL.Path[1:]
		w.Header().Set("Location", string(base58.Decode(shorturl)))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
