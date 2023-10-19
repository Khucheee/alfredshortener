package app

import (
	"github.com/btcsuite/btcutil/base58"
)

type URLData struct {
	originalurl string
	uuid        string
	isdeleted   bool
}

type Storage struct {
	Urls          map[string]URLData //мапа содержит сокращенный урл и полный
	keeper        Keeper
}

func NewStorage(keeper Keeper) *Storage {
	a := keeper.Restore()
	return &Storage{Urls: a, keeper: keeper}
}

func (s *Storage) AddURL(shorturl, url, uuid string) { //добавляем значение в мапу
	s.Urls[shorturl] = URLData{originalurl: url, uuid: uuid, isdeleted: false}
	s.keeper.Save(shorturl, url, uuid)
	//потом вызову save для keeper для сохранения на диск нового урла
}

func (s *Storage) CheckExistanse(originalurl string) string { //ищем значение в мапе, если "" то не нашли
	shorturl := base58.Encode([]byte(originalurl))
	url := s.Urls[shorturl].originalurl
	if url == "" {
		return ""
	}
	return shorturl
}

func (s *Storage) SearchURL(shorturl string) string { //ищем значение в мапе, если "" то не нашли
	url := s.Urls[shorturl].originalurl
	return url
}

func (s *Storage) getbyuser(uuid string) []Dburls {
	return s.keeper.GetUrlsByUser(uuid)
}

func (s *Storage) DeleteUserLinks(uid string, in []string) {
	//сначала удалю урл из памяти
	for _, hash := range in {
		structura := s.Urls[hash] //получаем из кэша данные по сокращенному урлу
		if structura.uuid!=uid {//если этот урл не принадлежит пользователю, берем следующее значение
			continue
		}
		structura.isdeleted = true //отмечаем в них что урл удален
		s.Urls[hash] = structura //записываем новые данные в кэш
	}

}
func (s *Storage) DeleteUserLink(uid,hash string) {
	s.keeper.DeleteUserLink(uid,hash )
	}

}
//либо в воркер прокидывать кипер
//