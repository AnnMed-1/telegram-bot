package main

import (
    "log"
    "os"
    tele "gopkg.in/tucnak/telebot.v2"
    "os/exec"
    "strings"
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

    log.Println("Бот запущен!")

    // Healthcheck endpoint
    go func() {
        http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
        })
        log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
    }()

    bot.Handle("/start", func(c *tele.Context) error {
        return c.Reply("🎉 Бот работает! Отправь ссылку на YouTube!")
    })

    bot.Handle(tele.OnText, func(c *tele.Context) error {
        url := c.Text()
        if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
            c.Notify(tele.Typing)
            err := downloadVideo(url, c.Chat())
            if err != nil {
                return c.Reply("❌ Ошибка: " + err.Error())
            }
            return c.Reply("✅ Видео скачано!")
        }
        return nil
    })

    bot.Start()
}

func downloadVideo(url string, chat tele.Recipient) error {
    // Создаём папку videos
    os.Mkdir("videos", 0755)

    // yt-dlp скачивает видео
    cmd := exec.Command("yt-dlp", 
        "-f", "best[height<=720]", 
        "--output", "videos/%(title)s.%(ext)s", 
        url)
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("yt-dlp failed: %s", string(output))
    }

    chat.Send(tele.Typing)
    return nil
}
