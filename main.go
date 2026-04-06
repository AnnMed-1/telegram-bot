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

    go func() {
        http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "text/plain")
            w.WriteHeader(http.StatusOK)
            fmt.Fprint(w, "OK")
        })
        log.Fatal(http.ListenAndServe(":"+port, nil))
    }()

    // Telegram Bot
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN required")
        return
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

    log.Println("✅ Skillspace Downloader запущен!")

    bot.Handle("/start", startHandler)
    bot.Handle("/help", helpHandler)
    bot.Handle(tele.OnText, skillspaceHandler)

    bot.Start()
}

func startHandler(c *tele.Context) error {
    msg := `📚 <b>Skillspace Video Downloader</b>

<b>Отправь ссылку Skillspace:</b>
• https://app.skillspace.ru/course/123/lesson/456
• https://app.skillspace.ru/courses/123/lessons/456

Автоматически найдёт MP4!`
    
    return c.Send(msg, &tele.SendOptions{ParseMode: tele.ModeHTML})
}

func helpHandler(c *tele.Context) error {
    return c.Reply("/start - Запуск\nОтправь Skillspace ссылку!")
}

func skillspaceHandler(c *tele.Context) error {
    url := strings.TrimSpace(c.Text())
    
    // Skillspace ссылки
    if !strings.Contains(url, "skillspace.ru") {
        return nil
    }

    c.Notify(tele.Typing)
    
    filename, err := downloadSkillspace(url)
    if err != nil {
        return c.Reply(fmt.Sprintf("❌ Ошибка:\n<code>%s</code>", err.Error()))
    }

    // Отправка
    videoFile := &tele.Document{File: tele.FromDisk(filename)}
    err = c.Send(videoFile)
    if err != nil {
        return c.Reply("✅ Скачано (>50MB не отправляется)")
    }

    return c.Reply("✅ Видео отправлено!")
}

func downloadSkillspace(url string) (string, error) {
    if err := os.MkdirAll("videos", 0755); err != nil {
        return "", fmt.Errorf("dir: %v", err)
    }

    // yt-dlp для Skillspace (работает отлично!)
    cmd := exec.Command("yt-dlp",
        "--cookies-from-browser", "chrome",  // Если нужны куки
        "-f", "best[ext=mp4]",
        "--output", "videos/%(title)s.%(ext)s",
        url,
    )

    log.Printf("Скачиваем Skillspace: %s", url)
    output, err := cmd.CombinedOutput()
    if err != nil {
        // Fallback без cookies
        cmd = exec.Command("yt-dlp",
            "-f", "best[ext=mp4]",
            "--output", "videos/%(title)s.%(ext)s",
            url,
        )
        output, err = cmd.CombinedOutput()
        if err != nil {
            return "", fmt.Errorf("yt-dlp: %s", string(output))
        }
    }

    files, err := filepath.Glob("videos/*.mp4")
    if err != nil || len(files) == 0 {
        return "", fmt.Errorf("mp4 not found")
    }

    log.Printf("✅ Skillspace: %s", files[0])
    return files[0], nil
}
