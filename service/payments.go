package service

import (
	"log"
	"os"
	"strconv"
	"work/logs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	token        = "6553780269:AAGKRvVeV7cswTqcjEErQKbBfdU6t6cYE-Y"
	paymentToken = os.Getenv("5420394252:TEST:543267")
)

func Pay(id int) {
	bot, _ := tgbotapi.NewBotAPI(token)
	invoice := tgbotapi.NewInvoice(
		int64(id),
		"Оплата за подписку!",
		"Платеж на сумму 4990₸",
		"custom_payload",
		paymentToken, // Токен для создания платежа
		"start_param",
		"KZT",
		&[]tgbotapi.LabeledPrice{{Label: "KZT", Amount: 499000}},
	)
	invoice.ProviderToken = "5420394252:TEST:543267"

	log.Println("Before sending invoice")
	_, _ = bot.Send(invoice)

	log.Println("After sending invoice")
}
func HandlePaymentRequest(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	logs.Log("@" + update.Message.From.UserName + "  " + "ИМЯ: " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "  " + "ID: " + strconv.Itoa(int(update.Message.Chat.ID)) + "  " + update.Message.Text + "\n")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `	Kaspi оплата 4990₸`)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Оплата", "https://pay.kaspi.kz/pay/jxrd4qnx"),
		),
	)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}
