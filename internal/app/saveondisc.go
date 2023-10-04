package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type JSONfile struct {
	UUID        string `json:"uuid"`
	Shorturl    string `json:"short_url"`
	Originalurl string `json:"original_url"`
}
type FileStorage struct {
	path    string
	counter int
}
type Keeper interface {
	Save(string, string, string) //теперь тут принимаем еще и uuid
	Restore() map[string]string
	GetUrlsByUser(string) []Dburls //дожны быть мтеоды save принимает мапу возвращает ошибку/restore не принимает ничего отдает мапу
}

/* надо создать файл
забрать его название из конфига
сохранить
*/

func NewKeeper(c Configure) *Keeper {
	var keeper Keeper
	if c.Dblink == "" {
		keeper = &FileStorage{path: c.FilePath, counter: 0}
		return &keeper
	}
	keeper = &Database{link: c.Dblink}

	return &keeper
}

func (k *FileStorage) Save(shorturl, originalurl, uuid string) {
	k.counter += 1
	//strconv.Itoa(k.counter) раньше uuid так объявлялся
	j := JSONfile{UUID: uuid, Shorturl: shorturl, Originalurl: originalurl} //собираем структуру для сборки jsonки
	data, err := json.Marshal(j)                                            //собираем из нее jsonку
	if err != nil {
		panic(err)
	}
	data = append(data, '\n')                                                   //добавляем разделитель
	file, err := os.OpenFile(k.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) //открываем файл
	if err != nil {
		fmt.Println(err)
	}
	_, err = file.Write(data) //засовываем в него нашу jsonку
	if err != nil {
		fmt.Println(err)
	}
	file.Close() //закрываем малыша
}

// этот метод не должен ничего принимать
// он должен вернуть мапу(или слайс когда то в будущем)
func (k *FileStorage) Restore() map[string]string {
	fmt.Println("сработал метод рестор файла")
	file, err := os.OpenFile(k.path, os.O_RDONLY|os.O_CREATE, 0666) //открываю файл
	if err != nil {
		fmt.Println(err)
	}
	text := bufio.NewReader(file) //сохраняю все его содержимое в переменную
	var data []byte
	urls := make(map[string]string)
	for { //создаю переменную для строк из файла
		data, err = text.ReadBytes('\n') //засовываю в переменную строку из файла
		if err != nil {
			fmt.Println("no data", err)
		}
		if len(data) == 0 { //если строка пустая, то перестаю читать
			break
		}
		jon := JSONfile{}                    //если строка не пустая, то создаю структуру под неё
		json.Unmarshal(data, &jon)           //парсю строку в эту структуру
		file.Close()                         // закрываю малыша
		urls[jon.Shorturl] = jon.Originalurl //добавляю в хранилище
	}
	return urls
}

func (k *FileStorage) GetUrlsByUser(uuid string) []Dburls {
	fmt.Println("Здесь еще нет реализации")
	return []Dburls{}
}
