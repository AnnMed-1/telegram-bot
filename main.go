package main

import (
    "log"
    "os"
    "os/exec"
    
    tele "gopkg.in/tucnak/telebot.v2"
)

func main() {
    bot, err := tele.NewBot(tele.Settings{
        Token: os.Getenv("BOT_TOKEN"),
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Bot started")

    bot.Handle("/start", func(c *tele.Context) error {
        return c.Reply("Отправь Skillspace ссылку!")
    })

    bot.Handle(tele.OnText, func(c *tele.Context) error {
        url := c.Text()
        c.Notify(tele.Typing)
        
        // yt-dlp скачивает Skillspace
        cmd := exec.Command("yt-dlp", url, "-o", "video.%(ext)s")
        _, err := cmd.Output()
        if err != nil {
            return c.Reply("Ошибка: " + err.Error())
        }
        
        video := &tele.Video{File: tele.FromDisk("video.mp4")}
        c.Send(video)
        
        return nil
    })

    bot.Start()
}
