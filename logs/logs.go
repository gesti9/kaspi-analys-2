package logs

import (
	"fmt"
	"os"
	"time"
)

func Log(n string) {
	// Получаем текущую дату и время
	currentTime := time.Now()

	// Форматируем дату и время в строку
	dateTimeString := currentTime.Format("2006-01-02 15:04:05")

	// Имя файла, в который мы будем записывать данные
	fileName := "logs/log.txt"

	// Открываем файл для записи (если файла нет, он будет создан)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Записываем дату и время в файл
	_, err = file.WriteString(dateTimeString + "    " + n)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}
}
