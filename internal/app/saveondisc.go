package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type JSONfile struct {
	UUID        string `json:"uuid"`
	Shorturl    string `json:"short_url"`
	Originalurl string `json:"original_url"`
}
type keeper struct {
	FileName string
}
type Keeper interface {
	Save(map[string]string, string, string)
	Restore() (string, string)
	//дожны быть мтеоды save принимает мапу возвращает ошибку/restore не принимает ничего отдает мапу
}

/* надо создать файл
забрать его название из конфига
сохранить
*/

func NewKeeper(filepath string) *keeper {
	return &keeper{FileName: filepath}
}

func (k *keeper) Save(foruuid map[string]string, shorturl string, originalurl string) {
	if k.FileName == "" { //если название файла пустое, то сразу отдыхаем
		return
	}
	j := JSONfile{UUID: strconv.Itoa(len(foruuid) + 1), Shorturl: shorturl, Originalurl: originalurl} //собираем структуру для сборки jsonки
	data, err := json.Marshal(j)                                                                      //собираем из нее jsonку
	if err != nil {
		panic(err)
	}
	data = append(data, '\n')                                                       //добавляем разделитель
	file, err := os.OpenFile(k.FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) //открываем файл
	if err != nil {
		fmt.Println(err)
	}
	_, err = file.Write(data) //засовываем в него нашу jsonку
	if err != nil {
		fmt.Println(err)
	}
	file.Close() //закрываем малыша
}

func (k *keeper) Restore() (string, string) {
	file, err := os.OpenFile(k.FileName, os.O_RDONLY|os.O_CREATE, 0666) //открываю файл
	if err != nil {
		fmt.Println(err)
	}
	text := bufio.NewReader(file)    //сохраняю все его содержимое в переменную
	var data []byte                  //создаю переменную для строк из файла
	data, err = text.ReadBytes('\n') //засовываю в переменную строку из файла
	if err != nil {
		fmt.Println("no data", err)
	}
	if len(data) == 0 { //если строка пустая, то перестаю читать
		return "", ""
	}
	jon := JSONfile{}          //если строка не пустая, то создаю структуру под неё
	json.Unmarshal(data, &jon) //парсю строку в эту структуру
	file.Close()               // закрываю малыша
	return jon.Shorturl, jon.Originalurl
}
