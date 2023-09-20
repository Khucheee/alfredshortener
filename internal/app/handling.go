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
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

//storage должен быть интерфейсом!
//сторадж должен просто сохранять данные
//в этом файле описываю storage
//интерфейс описывается там где используется!!!

// теперь каким будет код: проверки на базу не нужны, проверки на пустоту файла внутри метода, то есть просто вызываем методы
func NewBaseController(c Configure, s Storage, l Logger) *BaseController {
	return &BaseController{config: c, storage: s, logger: l}
}

func (b *BaseController) solvePost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)

	if shorturl := b.storage.CheckExistanse(string(reqBody)); shorturl != "" {
		respBody := b.config.Address + shorturl
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(respBody))
		return
	}

	w.Header().Set("Content-Type", "text/plain")      //установили заголовок
	w.WriteHeader(http.StatusCreated)                 //установили статускод
	reqBodyEncoded := base58.Encode(reqBody)          //закодировали ссылку
	defer r.Body.Close()                              //закрыли тело
	respBody := b.config.Address + reqBodyEncoded     //собрали тело ответа
	w.Write([]byte(respBody))                         //отправили его
	b.storage.AddURL(reqBodyEncoded, string(reqBody)) //сохранили ссылку в мапу, а внутри дальше пойдет в сторадж
}

func (b *BaseController) solveGet(w http.ResponseWriter, r *http.Request) {
	if b.storage.SearchURL(chi.URLParam(r, "shorturl")) == "" { //если ключ в мапе пустой, то 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", b.storage.Urls[chi.URLParam(r, "shorturl")]) //если дошли до сюда, то в location суем значение из мапы по ключу
	w.WriteHeader(http.StatusTemporaryRedirect)
	return
}

func (b *BaseController) solveJSON(w http.ResponseWriter, r *http.Request) {
	var jsonquery Jsonquery
	var jsonresponse Jsonresponse
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body) //читаем тело запроса в буфер
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &jsonquery); err != nil { //парсим полученное счастье в структуру для запроса
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json") //устанавливаем заголовок

	if shorturl := b.storage.CheckExistanse(jsonquery.URL); shorturl != "" {
		jsonresponse.Response = b.config.Address + shorturl
		resp, _ := json.Marshal(jsonresponse)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(resp))
		return
	}

	shorturl := base58.Encode([]byte(jsonquery.URL))    //кодируем урл
	jsonresponse.Response = b.config.Address + shorturl //собираем тело ответа
	w.WriteHeader(http.StatusCreated)                   //устанавливаем статус код 201
	resp, _ := json.Marshal(jsonresponse)               //парсим его в json
	w.Write(resp)                                       //отправляем
	b.storage.AddURL(shorturl, jsonquery.URL)           //сохраняем в мапу, а она внутри сохранит еще куда надо

}

func (b *BaseController) solvePing(w http.ResponseWriter, r *http.Request) {
	if DBconnect(&b.config) {
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
		b.storage.AddURL(shorturl, request.OriginalURL)
	}

	resp, _ := json.Marshal(surls)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
