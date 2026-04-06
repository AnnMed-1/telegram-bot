package main

import (
    "log"
    "os"
    "fmt"
    "net/http"
    "os/exec"
    "strings"
    "time"
    
    tele "gopkg.in/tucnak/telebot.v2"
)

var bot *tele.Bot

func main() {
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN не найден! Добавь в Railway Variables")
    }

    pref := tele.Settings{
        Token:  token,
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    }

    var err error
    bot, err = tele.NewBot(pref)
    if err != nil {
        log.Fatal(err)
        return
    }

    log.Println("✅ Бот запущен!")

    // Healthcheck для Railway
    go func() {
        http.HandleFunc("/health", func(w http.ResponseWriter, r *tele.Request) {
            w.WriteHeader(http.StatusOK)
        })
        port := os.Getenv("PORT")
        if port == "" {
            port = "8080"
        }
        log.Fatal(http.ListenAndServe(":"+port, nil))
    }()

    bot.Handle("/start", func(c *tele.Context) error {
        return c.Reply("🎉 Бот работает! Отправь ссылку YouTube!")
    })

    bot.Handle(tele.OnText, func(c *tele.Context) error {
        url := c.Text()
        if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
            c.Notify(tele.Typing)
            err := downloadVideo(url, c.Chat())
            if err != nil {
                return c.Reply("❌ Ошибка: " + err.Error())
            }
            return c.Reply("✅ Видео готово для скачивания!")
        }
        return nil
    })

    bot.Start()
}

func downloadVideo(url string, chat tele.Recipient) error {
    os.Mkdir("videos", 0755)

    cmd := exec.Command("yt-dlp", 
        "-f", "best[height<=720]", 
        "--output", "videos/%(title)s.%(ext)s", 
        url)
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("yt-dlp failed: %s", string(output))
    }

    return nil
}
