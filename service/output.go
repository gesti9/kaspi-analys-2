package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func Output(n string) (string, error) {
	url := n
	var result string

	// Отправка GET-запроса
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	// Выполнение запроса с http.DefaultClient
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ошибка запроса. Статус: %d", resp.StatusCode)
	}

	// Используем goquery для парсинга HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании документа из ответа: %v", err)
	}

	// Используем регулярное выражение для поиска значения "reviewsCount"
	re := regexp.MustCompile(`"reviewsCount":(\d+)`)

	// Поиск по тексту страницы
	doc.Find("script").Each(func(index int, element *goquery.Selection) {
		scriptText := element.Text()
		match := re.FindStringSubmatch(scriptText)
		if len(match) == 2 {
			reviewsCount := match[1]
			fmt.Printf("Значение reviewsCount: %s\n", reviewsCount)

			result = reviewsCount
		}
	})

	return strconv.Itoa(int(Sum(result))), nil
}

func Price(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("ошибка при выполнении GET-запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка при чтении тела ответа: %v", err)
	}

	// Найти все совпадения с регулярным выражением
	re := regexp.MustCompile(`"price": "(\d+)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	// Извлечь только цифры из совпадений
	var sum int
	for _, match := range matches {
		if len(match) >= 2 {
			price, err := strconv.Atoi(match[1])
			if err != nil {
				return 0, fmt.Errorf("ошибка при конвертации строки в число: %v", err)
			}
			sum += price
		}
	}

	return sum, nil
}

// Функция для вычисления суммы (ваша логика)
func Sum(s string) float64 {
	num, _ := strconv.Atoi(s)
	return float64(57*num) / float64(43)
}
