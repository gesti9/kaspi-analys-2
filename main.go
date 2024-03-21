package main

import (
	"fmt"
	"log"
	"strconv"
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
			count, _ := strconv.Atoi(data.ReadFromFile("data/users/" + strconv.Itoa(int(update.Message.Chat.ID)) + ".txt"))
			switch update.Message.Text {
			case "/start":
				logs.Log("@" + update.Message.From.UserName + "  " + "ИМЯ: " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "  " + "ID: " + strconv.Itoa(int(update.Message.Chat.ID)) + "  " + update.Message.Text + "\n")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для анализа отправьте ссылку с Kaspi.kz")
				msg.ReplyToMessageID = update.Message.MessageID
				if update.Message != nil && update.Message.Contact != nil {
					// Теперь можно безопасно использовать PhoneNumber
					fmt.Println(update.Message.Contact.PhoneNumber)
				} else {
					fmt.Println("Контакт не предоставлен или объект Contact равен nil")
				}
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)

			case "Администратор!":
				logs.Log("@" + update.Message.From.UserName + "  " + "ИМЯ: " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "  " + "ID: " + strconv.Itoa(int(update.Message.Chat.ID)) + "  " + update.Message.Text + "\n")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Админ🐼`)
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL("Админ тут", "https://t.me/Alisher_Adenov"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			case "/admin":
				logs.Log("@" + update.Message.From.UserName + "  " + "ИМЯ: " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "  " + "ID: " + strconv.Itoa(int(update.Message.Chat.ID)) + "  " + update.Message.Text + "\n")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Переходи👇`)
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL("Поддержка!", "https://t.me/Alisher_Adenov"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			default:
				logs.Log("@" + update.Message.From.UserName + "  " + "ИМЯ: " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "  " + "ID: " + strconv.Itoa(int(update.Message.Chat.ID)) + "  " + update.Message.Text + "\n")
				if service.IsValidURL(update.Message.Text) {
					fmt.Printf("%s - это валидная ссылка\n", (update.Message.Text))
					result, _ := service.Output(update.Message.Text)
					num, _ := strconv.Atoi(result)

					if data.ReadFromFile("data/users/"+strconv.Itoa(int(update.Message.Chat.ID))+".txt") == "10" {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для продолжения оплатите 4000 тенге, для оплаты напишите Администратору!")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
						fmt.Println(data.ReadFromFile("data/users/" + strconv.Itoa(int(update.Message.Chat.ID)) + ".txt"))
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
						fmt.Println(formatted)
						price, _ := service.Price(update.Message.Text)
						moneyM := price * 30
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, `Продажи за все время: `+result+` шт`+"\n"+
							`За месяц: `+strconv.Itoa(int(mes))+` шт`+"\n"+`За день: `+formatted+` шт`+
							"\n"+`Заработали: `+strconv.Itoa(moneyM)+` тенге`)
						msg.ReplyToMessageID = update.Message.MessageID

						bot.Send(msg)

						count++
						data.UserData(update.Message.From.ID, count)
					}

					// Ваш код для загрузки и анализа страницы
				} else {
					fmt.Printf("%s - не является валидной ссылкой\n", (update.Message.Text))
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для анализа отправьте ссылку с Kaspi.kz")
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
				}

			}

		}
	}
}
