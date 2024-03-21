package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func UserData(chat_id int64, s int) {
	// Получаем текущую дату и время

	// Имя файла, в который мы будем записывать данные
	fileName := "data/users/" + strconv.Itoa(int(chat_id)) + ".txt"

	// Открываем файл для записи (если файла нет, он будет создан)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Записываем дату и время в файл
	_, err = file.WriteString(strconv.Itoa(s))
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}
	// fmt.Println(ReadFromFile(fileName))
	if ReadFromFile(fileName) == "2" {
		fmt.Println("kuku")
	}
}

func ReadFromFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ""
	}
	return string(content)
}
