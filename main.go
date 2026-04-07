package main

import (
    "log"
    "os"
    "time"
    
    tele "gopkg.in/tucnak/telebot.v2"
)

func main() {
    token := os.Getenv("BOT_TOKEN")
    log.Printf("Запуск с токеном: %s...", token[:10])
    
    bot, err := tele.NewBot(tele.Settings{
        Token:  token,
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal("BOT ERROR:", err)
        return
    }

    bot.Handle("/start", func(m *tele.Message) {
        bot.Send(m.Sender, "✅ @Kurses_skil_bot готов!")
    })

    log.Println("✅ Bot started!")
    bot.Start()
}
