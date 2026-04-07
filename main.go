package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    
    tele "gopkg.in/telebot.v2"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" { port = "8080" }
    
    // Healthcheck
    go func() {
        http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "text/plain")
            w.WriteHeader(200)
            fmt.Fprint(w, "ok")
        })
        log.Printf("✅ Healthcheck :%s/health", port)
        log.Fatal(http.ListenAndServe(":"+port, nil))
    }()

    // Bot LongPoller
    token := os.Getenv("BOT_TOKEN")
    bot, err := tele.NewBot(tele.Settings{
        Token:  token,
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal("❌ BOT ERROR:", err)
        return
    }

    bot.Handle("/start", func(c tele.Context) error {
        return c.Reply("✅ @Kurses_skil_bot готов!")
    })

    log.Println("✅ Bot started - пиши /start!")
    bot.Start()
}
