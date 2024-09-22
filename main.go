package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"work/data"
	"work/logs"
	"work/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserState struct {
	CurrentState string
	PrevState    string
}

var (
	bot             *tgbotapi.BotAPI
	userStates      = make(map[int64]*UserState)
	userStatesMutex sync.Mutex
	mainMenu        = tgbotapi.NewReplyKeyboard(

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Администратор!"),
		),
	)
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6831494639:AAGkAcG9BgZYarNfcviU-SsH3hvnadcLjkE")
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Читаем текущее значение из файла
			countStr := data.ReadFromFile("data/users/" + strconv.Itoa(int(update.Message.Chat.ID)) + ".txt")
			countStr = strings.TrimSpace(countStr) // Удаляем пробелы и символы новой строки
			count, _ := strconv.Atoi(countStr)

			switch update.Message.Text {
			case "/start":
				logs.Log("@" + update.Message.From.UserName + "  " + "ИМЯ: " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "  " + "ID: " + strconv.Itoa(int(update.Message.Chat.ID)) + "  " + update.Message.Text + "\n")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для анализа отправьте ссылку с Kaspi.kz")
				msg.ReplyToMessageID = update.Message.MessageID
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)

			case "Администратор!", "/admin":
				logs.Log("@" + update.Message.From.UserName + "  " + "ИМЯ: " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "  " + "ID: " + strconv.Itoa(int(update.Message.Chat.ID)) + "  " + update.Message.Text + "\n")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Админ🐼`)
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL("Админ тут", "https://t.me/dba_nurs"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			default:
				logs.Log("@" + update.Message.From.UserName + "  " + "ИМЯ: " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "  " + "ID: " + strconv.Itoa(int(update.Message.Chat.ID)) + "  " + update.Message.Text + "\n")
				if service.IsValidURL(update.Message.Text) {
					fmt.Printf("%s - это валидная ссылка\n", update.Message.Text)
					result, _ := service.Output(update.Message.Text)
					num, _ := strconv.Atoi(result)

					if count == 10000 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для продолжения оплатите 500 тенге для доступа на неограниченное количество запросов на 1 месяц, для оплаты напишите Администратору!")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					} else if num == 0 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "0 продаж!")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваш запрос обрабатывается..")
						bot.Send(msg)

						mes := (float64(num) / float64(365)) * 30
						day := float64(mes) / float64(30)
						formatted := fmt.Sprintf("%.1f", day)
						price, _ := service.Price(update.Message.Text)
						var moneyM, moneyM2 int
						res1, _ := strconv.Atoi(result)

						if int(mes) == 0 {
							moneyM = price * res1
							moneyM2 = 0
						} else {
							moneyM = price * res1
							moneyM2 = price * int(mes)
						}

						msg = tgbotapi.NewMessage(update.Message.Chat.ID, `Продажи за все время: `+result+` шт`+"\n"+
							`За месяц: `+strconv.Itoa(int(mes))+` шт`+"\n"+`За день: `+formatted+` шт`+"\n"+`В месяц заработок: `+strconv.Itoa(moneyM2)+` тенге`+
							"\n"+`Заработали за все время: `+strconv.Itoa(moneyM)+` тенге`)
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)

						count++
						data.UserData(update.Message.From.ID, count)
					}

				} else {
					fmt.Printf("%s - не является валидной ссылкой\n", update.Message.Text)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для анализа отправьте ссылку с Kaspi.kz")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
			}
		}
	}

}
