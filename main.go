package main

import (
    "log"
    "os"
    "time"
    
    tele "gopkg.in/tucnak/telebot.v2"
)

func main() {
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN пустой!")
        return
    }
    
    bot, err := tele.NewBot(tele.Settings{
        Token:  token,
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal("BOT ERROR:", err)
        return
    }

    bot.Handle("/start", func(c tele.Context) error {
        return c.Reply("✅ Skillspace Downloader готов!")
    })

    log.Println("✅ Bot started!")
    bot.Start()
}
