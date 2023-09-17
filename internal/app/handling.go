package app

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// структуры для json хендлера
type Jsonquery struct {
	URL string
}
type Jsonresponse struct {
	Response string `json:"result"`
}

// структуры для ручки batch
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func NewBaseController(c Configure, s Storage, l Logger) *BaseController {
	return &BaseController{config: c, storage: s, logger: l}
}

func (b *BaseController) solvePost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	if shorturl := GetShortUrldb(string(reqBody), b.config); shorturl != "" {
		respBody := b.config.Address + shorturl
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(respBody))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	reqBodyEncoded := base58.Encode(reqBody)

	defer r.Body.Close()
	respBody := b.config.Address + reqBodyEncoded
	w.Write([]byte(respBody))
	if b.config.Dblink == "" {
		b.storage.keeper.Save(b.storage.Urls, reqBodyEncoded, string(reqBody))
		b.storage.AddURL(reqBodyEncoded, string(reqBody))
		return
	}
	err := AddURLdb(reqBodyEncoded, string(reqBody), b.config)
	if err != nil {
		b.logger.sugar.Infoln(err)
	}
}

func (b *BaseController) solveGet(w http.ResponseWriter, r *http.Request) {
	if b.config.Dblink == "" {
		if b.storage.SearchURL(chi.URLParam(r, "shorturl")) == "" { //если ключ в мапе пустой, то 400
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", b.storage.Urls[chi.URLParam(r, "shorturl")]) //если дошли до сюда, то в location суем значение из мапы по ключу
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if GetOriginalUrldb(chi.URLParam(r, "shorturl"), b.config) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", GetOriginalUrldb(chi.URLParam(r, "shorturl"), b.config)) //если дошли до сюда, то в location суем значение из мапы по ключу
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (b *BaseController) solveJSON(w http.ResponseWriter, r *http.Request) {
	var jsonquery Jsonquery
	var jsonresponse Jsonresponse
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &jsonquery); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if shorturl := GetShortUrldb(jsonquery.URL, b.config); shorturl != "" {
		jsonresponse.Response = b.config.Address + shorturl
		fmt.Println(shorturl)
		resp, _ := json.Marshal(jsonresponse)
		w.WriteHeader(http.StatusConflict)
		w.Write(resp)
		return
	}
	shorturl := base58.Encode([]byte(jsonquery.URL))
	jsonresponse.Response = b.config.Address + shorturl
	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(jsonresponse)
	w.Write(resp)

	if b.config.Dblink == "" {
		b.storage.AddURL(shorturl, jsonquery.URL)
		b.storage.keeper.Save(b.storage.Urls, shorturl, jsonquery.URL)
		return
	}
	AddURLdb(shorturl, jsonquery.URL, b.config)
}

func (b *BaseController) solvePing(w http.ResponseWriter, r *http.Request) {
	if DBconnect(b.config) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (b *BaseController) solveBatch(w http.ResponseWriter, r *http.Request) {
	var ourls = []BatchRequest{}
	var surls = []BatchResponse{}
	var response BatchResponse
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &ourls); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, request := range ourls {
		shorturl := base58.Encode([]byte((request.OriginalURL)))
		response.CorrelationID = request.CorrelationID
		response.ShortURL = b.config.Address + shorturl
		surls = append(surls, response)
		if b.config.Dblink == "" {
			b.storage.keeper.Save(b.storage.Urls, shorturl, request.OriginalURL)
			b.storage.AddURL(shorturl, request.OriginalURL)
			continue
		}
		AddURLdb(shorturl, request.OriginalURL, b.config)
	}
	resp, _ := json.Marshal(surls)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
