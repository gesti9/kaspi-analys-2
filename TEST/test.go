package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

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

func main() {
	url := "https://kaspi.kz/shop/p/elektrochainik-bereke-br-810-seryi-109981423/?c=710000000"
	price, _ := Price(url)
	fmt.Println(price + 1000)
}
