package app

type Storage struct {
	Urls   map[string]string //мапа содержит сокращенный урл и полный
	keeper Keeper
}

func NewStorage(keeper Keeper) *Storage {
	return &Storage{Urls: make(map[string]string), keeper: keeper}
}

func (s *Storage) AddURL(shorturl, url string) { //добавляем значение в мапу
	s.Urls[shorturl] = url
	//потом вызову save для keeper для сохранения на диск нового урла
}

func (s *Storage) SearchURL(shorturl string) string { //ищем значение в мапе, если "" то не нашли
	url := s.Urls[shorturl]
	return url
}
