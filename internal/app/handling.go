package app

import (
	"github.com/btcsuite/btcutil/base58"
	"io"
	"net/http"
)

var Shorturl string

func SalamPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusCreated)
		Body, _ := io.ReadAll(r.Body)
		Shorturl = base58.Encode(Body)
		w.Write([]byte("http://localhost:8080/" + Shorturl))
	} else if r.Method == http.MethodGet && r.URL.Path == "/"+Shorturl {
		w.Header().Set("location", string(base58.Decode(Shorturl)))
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
