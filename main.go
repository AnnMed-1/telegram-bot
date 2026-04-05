package main

import (
    "log"
    "os"
    "strings"
    "time"
    
    "gopkg.in/telebot.v2"
)

func main() {
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN is empty")
    }

    bot, err := telebot.NewBot(telebot.Settings{
        Token:  token,
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println("🎥 SkillSpace Downloader запущен!")

    bot.Handle("/start", func(m *telebot.Message) {
        bot.Send(m.Sender, 
`🎥 SkillSpace Video Downloader

📱 Отправь ссылку на видео SkillSpace:
https://skillspace.ru/embed/...
https://kinescope.io/...

⏳ Скачаю Full HD MP4 за 1-2 мин!`)
    })

    bot.Handle(telebot.OnText, func(m *telebot.Message) {
        url := strings.TrimSpace(m.Text)
        
        if strings.Contains(url, "skillspace") || strings.Contains(url, "kinescope") {
            bot.Send(m.Sender, "⏳ Ищу SkillSpace видео...")
            
            // Пока заглушка — потом yt-dlp
            bot.Send(m.Sender, 
`✅ Видео готово! (демо)

📹 Full HD MP4 скачан
⏱️ Время: 1:23
💾 Размер: 245MB

[yt-dlp работает в фоне]`)
        } else {
            bot.Send(m.Sender, "❌ Только SkillSpace/Kinescope ссылки!")
        }
    })

    bot.Start()
}
// fix token cache - 2026-04-05
