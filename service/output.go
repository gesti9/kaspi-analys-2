package service

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
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
	// Run Chrome browser
	service, err := selenium.NewChromeDriverService("C:/chromedriver-win64/chromedriver.exe", 4444)
	if err != nil {
		return 0, err
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"disable-gpu",
		"--headless", // раскомментируйте эту строку, чтобы сделать браузер невидимым
	}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		return 0, err
	}

	driver.Get(url)

	// Ждем 2 секунды вместо 5
	time.Sleep(2 * time.Second)

	// Находим элемент по классу
	elem, err := driver.FindElement(selenium.ByClassName, "item__price-once")
	if err != nil {
		return 0, err
	}

	// Получаем текст из элемента
	text, err := elem.Text()
	if err != nil {
		return 0, err
	}

	// Извлекаем только цифры с использованием регулярного выражения
	re := regexp.MustCompile(`\d+`)
	digits := re.FindAllString(text, -1)

	// Объединяем извлеченные цифры в одну строку
	resultString := ""
	for _, digit := range digits {
		resultString += digit
	}

	// Преобразуем строку в число
	res, err := strconv.Atoi(resultString)
	if err != nil {
		return 0, err
	}

	fmt.Println("Текст из элемента item__price-once:", text)
	fmt.Println("Извлеченные цифры:", res)

	return res, nil
}

// Функция для вычисления суммы (ваша логика)
func Sum(s string) float64 {
	num, _ := strconv.Atoi(s)
	return float64(57*num) / float64(43)
}
