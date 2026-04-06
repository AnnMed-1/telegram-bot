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
    // Healthcheck
    go func() {
        port := os.Getenv("PORT")
        if port == "" { port = "8080" }
        http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "text/plain")
            w.WriteHeader(http.StatusOK)
            fmt.Fprint(w, "OK")
        })
        log.Fatal(http.ListenAndServe(":"+port, nil))
    }()

    // Bot
    token := os.Getenv("BOT_TOKEN")
    bot, err := tele.NewBot(tele.Settings{
        Token:  token,
        Poller: &tele.LongPoller{},
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println("✅ Bot started!")

    bot.Handle("/start", func(c *tele.Context) error {
        return c.Reply("📚 Skillspace готов!")
    })

    bot.Handle(tele.OnText, func(c *tele.Context) error {
        c.Notify(tele.Typing)
        pattern := filepath.Join(os.TempDir(), "video.%(ext)s")
        cmd := exec.Command("yt-dlp", "-o", pattern, c.Text())
        if out, err := cmd.CombinedOutput(); err != nil {
            return c.Reply("❌ " + string(out))
        }
        files, _ := filepath.Glob(filepath.Join(os.TempDir(), "video.*"))
        if len(files) > 0 {
            c.Send(&tele.Video{File: tele.FromDisk(files[0])})
            os.Remove(files[0])
        }
        return nil
    })

    bot.Start()
}
