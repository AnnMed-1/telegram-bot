package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
    
    tele "gopkg.in/tucnak/telebot.v2"
)

var bot *tele.Bot

func main() {
    // Railway PORT + Healthcheck
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Healthcheck endpoint
    go func() {
        http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "text/plain")
            w.WriteHeader(http.StatusOK)
            fmt.Fprint(w, "OK")
        })
        log.Printf("Healthcheck on :%s/health", port)
        log.Fatal(http.ListenAndServe(":"+port, nil))
    }()

    // Telegram Bot
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN not found! Add to Railway Variables")
        return
    }

    pref := tele.Settings{
        Token:  token,
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    }

    var err error
    bot, err = tele.NewBot(pref)
    if err != nil {
        log.Fatal("Bot error:", err)
        return
    }

    log.Println("✅ Telegram Bot запущен!")

    // Команды
    bot.Handle("/start", startHandler)
    bot.Handle("/help", helpHandler)
    bot.Handle(tele.OnText, videoHandler)

    // Запуск
    bot.Start()
}

func startHandler(c *tele.Context) error {
    msg := `🎬 <b>YouTube Downloader Bot</b>

<b>Отправь ссылку на YouTube:</b>
• https://youtube.com/watch?v=abc123
• https://youtu.be/abc123

<b>Поддерживает:</b> MP4 720p max`
    
    return c.Send(msg, &tele.SendOptions{ParseMode: tele.ModeHTML})
}

func helpHandler(c *tele.Context) error {
    return c.Reply(`🤖 Просто отправь YouTube ссылку!

/start - Запуск
/help - Помощь`)
}

func videoHandler(c *tele.Context) error {
    url := strings.TrimSpace(c.Text())
    
    // Проверка YouTube
    if !strings.Contains(url, "youtube.com") && !strings.Contains(url, "youtu.be") {
        return nil
    }

    // Typing indicator
    c.Notify(tele.Typing)

    // Скачивание
    filename, err := downloadVideo(url)
    if err != nil {
        return c.Reply(fmt.Sprintf("❌ Ошибка скачивания:\n<code>%s</code>", err.Error()))
    }

    // Отправка видео
    videoFile := &tele.Document{File: tele.FromDisk(filename)}
    err = c.Send(videoFile)
    if err != nil {
        return c.Reply("✅ Скачано, но не могу отправить (>50MB)")
    }

    return c.Reply("✅ Видео отправлено!")
}

func downloadVideo(url string) (string, error) {
    // Создаём папку
    if err := os.MkdirAll("videos", 0755); err != nil {
        return "", fmt.Errorf("create dir: %v", err)
    }

    // yt-dlp команда
    cmd := exec.Command("yt-dlp",
        "-f", "best[height<=720]",
        "--merge-output-format", "mp4",
        "-o", "videos/%(title)s.%(ext)s",
        url,
    )

    log.Printf("Скачиваем: %s", url)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("yt-dlp: %s", string(output))
    }

    // Найти файл
    files, err := filepath.Glob("videos/*.mp4")
    if err != nil || len(files) == 0 {
        return "", fmt.Errorf("no video found")
    }

    log.Printf("Скачано: %s (%d bytes)", files[0], getFileSize(files[0]))
    return files[0], nil
}

func getFileSize(path string) int64 {
    info, err := os.Stat(path)
    if err != nil {
        return 0
    }
    return info.Size()
}
