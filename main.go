package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    
    tele "gopkg.in/tucnak/telebot.v2"
)

func main() {
    // Healthcheck для Railway
    go func() {
        http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(200)
        })
        port := os.Getenv("PORT")
        if port == "" { port = "8080" }
        log.Fatal(http.ListenAndServe(":"+port, nil))
    }()

    // Bot с LongPoller (БЕЗ webhook!)
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN required")
    }

    bot, err := tele.NewBot(tele.Settings{
        Token:  token,
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println("✅ Bot started!")

    // ПРАВИЛЬНЫЕ handlers
    bot.Handle("/start", func(c *tele.Context) error {
        return c.Reply("📚 Skillspace Downloader!\nОтправь ссылку app.skillspace.ru")
    })

    bot.Handle(tele.OnText, func(c *tele.Context) error {
        url := c.Text()
        if !strings.Contains(url, "skillspace.ru") {
            return nil
        }

        c.Notify(tele.Typing)
        
        // yt-dlp в temp
        tmpDir := os.TempDir()
        pattern := filepath.Join(tmpDir, "skillspace.%(ext)s")
        
        cmd := exec.Command("yt-dlp", "-f", "best[ext=mp4]", "-o", pattern, url)
        out, err := cmd.CombinedOutput()
        if err != nil {
            return c.Reply("❌ " + string(out))
        }

        // Найти файл
        files, _ := filepath.Glob(filepath.Join(tmpDir, "skillspace.*"))
        if len(files) == 0 {
            return c.Reply("Файл не найден")
        }

        // Отправить
        video := &tele.Video{File: tele.FromDisk(files[0])}
        if err := c.Send(video); err != nil {
            c.Reply("✅ Скачано, но >50MB")
        }

        os.Remove(files[0]) // Cleanup
        return nil
    })

    bot.Start()
}
