package app

type Storage struct {
	Urls map[string]string //мапа содержит сокращенный урл и полный
	//тут будет keeper
}

type Keeper interface {
	AddURL(shorturl, url string)
	SearchURL(shorturl string) string
	//дожны быть мтеоды save принимает мапу возвращает ошибку/restore не принимает ничего отдает мапу
}

func (s *Storage) AddURL(shorturl, url string) { //добавляем значение в мапу
	s.Urls[shorturl] = url
	//потом вызову save для keeper для сохранения на диск нового урла
}

func (s *Storage) SearchURL(shorturl string) string { //ищем значение в мапе, если "" то не нашли
	url := s.Urls[shorturl]
	return url
}
