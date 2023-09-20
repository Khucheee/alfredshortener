package app

import "github.com/btcsuite/btcutil/base58"

type Storage struct {
	Urls   map[string]string //мапа содержит сокращенный урл и полный
	keeper Keeper
}

func NewStorage(keeper Keeper) *Storage {
	return &Storage{Urls: make(map[string]string), keeper: keeper}
}

func (s *Storage) AddURL(shorturl, url string) { //добавляем значение в мапу
	s.Urls[shorturl] = url
	s.keeper.Save(shorturl, url)
	//потом вызову save для keeper для сохранения на диск нового урла
}
func (s *Storage) CheckExistanse(originalurl string) string { //ищем значение в мапе, если "" то не нашли
	shorturl := base58.Encode([]byte(originalurl))
	url := s.Urls[shorturl]
	if url == "" {
		return ""
	}
	return shorturl
}

func (s *Storage) SearchURL(shorturl string) string { //ищем значение в мапе, если "" то не нашли
	url := s.Urls[shorturl]
	return url
}

func (s *Storage) Restore() {
	s.Urls = s.keeper.Restore()
}
