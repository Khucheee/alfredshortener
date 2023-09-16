package app

import (
	"bytes"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

// контроллер для хендлеров
type BaseController struct {
	config  Configure
	storage Storage
	logger  Logger
}

type Jsonquery struct {
	URL string
}
type Jsonresponse struct {
	Response string `json:"result"`
}

func NewBaseController(c Configure, s Storage, l Logger) *BaseController {
	return &BaseController{config: c, storage: s, logger: l}
}

func (b *BaseController) solvePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	reqBody, _ := io.ReadAll(r.Body)
	reqBodyEncoded := base58.Encode(reqBody)

	defer r.Body.Close()
	respBody := b.config.Address + reqBodyEncoded
	w.Write([]byte(respBody))
	b.storage.keeper.Save(b.storage.Urls, reqBodyEncoded, string(reqBody))
	b.storage.AddURL(reqBodyEncoded, string(reqBody))

}

func (b *BaseController) solveGet(w http.ResponseWriter, r *http.Request) {
	if b.storage.SearchURL(chi.URLParam(r, "shorturl")) == "" { //если ключ в мапе пустой, то 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", b.storage.Urls[chi.URLParam(r, "shorturl")]) //если дошли до сюда, то в location суем значение из мапы по ключу
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (b *BaseController) solveJSON(w http.ResponseWriter, r *http.Request) {
	var jsonquery Jsonquery
	var jsonresponse Jsonresponse
	var buf bytes.Buffer
	var shorturl string
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &jsonquery); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shorturl = base58.Encode([]byte(jsonquery.URL))
	jsonresponse.Response = b.config.Address + shorturl
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(jsonresponse)
	w.Write(resp)
	b.storage.AddURL(shorturl, jsonquery.URL)
	b.storage.keeper.Save(b.storage.Urls, shorturl, jsonquery.URL)
}

func (b *BaseController) solvePing(w http.ResponseWriter, r *http.Request) {
	if DBconnect(b.config) == true {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
