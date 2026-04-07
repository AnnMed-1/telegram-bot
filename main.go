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
        Poller
