package main

import (
	"log"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	token := "8703503369:AAEDz5vboyJ7z9wizkwxQOdANgnGg5opirY"  // Твой токен!
	
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(m *tb.Message) {
		bot.Send(m.Sender, "Привет! Бот работает! 🚀")
	})

	log.Println("Бот запущен!")
	bot.Start()
}