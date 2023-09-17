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

// структуры для json хендлера
type Jsonquery struct {
	URL string
}
type Jsonresponse struct {
	Response string `json:"result"`
}

// структуры для ручки batch
type BatchRequest struct {
	Correlation_id string `json:"correlation_id"`
	Original_url   string `json:"original_url"`
}

type BatchResponse struct {
	Correlation_id string `json:"correlation_id"`
	Short_url      string `json:"short_url"`
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
	if b.config.Dblink == "" {
		b.storage.keeper.Save(b.storage.Urls, reqBodyEncoded, string(reqBody))
		b.storage.AddURL(reqBodyEncoded, string(reqBody))
		return
	}
	AddURLdb(reqBodyEncoded, string(reqBody), b.config)
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
	if GetUrldb(chi.URLParam(r, "shorturl"), b.config) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", GetUrldb(chi.URLParam(r, "shorturl"), b.config)) //если дошли до сюда, то в location суем значение из мапы по ключу
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
		shorturl := base58.Encode([]byte((request.Original_url)))
		response.Correlation_id = request.Correlation_id
		response.Short_url = shorturl
		surls = append(surls, response)
		if b.config.Dblink == "" {
			b.storage.AddURL(shorturl, request.Original_url)
			b.storage.keeper.Save(b.storage.Urls, shorturl, request.Original_url)
			continue
		}
		AddURLdb(shorturl, request.Original_url, b.config)
	}
	resp, _ := json.Marshal(surls)
	w.Write(resp)
}
