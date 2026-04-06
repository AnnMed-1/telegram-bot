package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	tele "gopkg.in/tucnak/telebot.v2"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is empty")
	}

	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		log.Fatal("PUBLIC_URL is empty")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	webhook := &tele.Webhook{
		Listen: ":" + port,
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: publicURL,
		},
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: webhook,
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(m *tele.Message) {
		_, _ = bot.Send(m.Sender, "Отправь ссылку")
	})

	bot.Handle(tele.OnText, func(m *tele.Message) {
		url := m.Text

		tmpDir := os.TempDir()
		pattern := filepath.Join(tmpDir, "video.%(ext)s")

		cmd := exec.Command("yt-dlp", "-o", pattern, url)
		out, err := cmd.CombinedOutput()
		if err != nil {
			_, _ = bot.Send(m.Sender, "Ошибка:\n"+string(out))
			return
		}

		files, err := filepath.Glob(filepath.Join(tmpDir, "video.*"))
		if err != nil || len(files) == 0 {
			_, _ = bot.Send(m.Sender, "Файл скачан, но не найден")
			return
		}

		video := &tele.Video{File: tele.FromDisk(files[0])}
		_, err = bot.Send(m.Sender, video)
		if err != nil {
			_, _ = bot.Send(m.Sender, "Не удалось отправить видео: "+err.Error())
		}

		_ = os.Remove(files[0])
	})

	go func() {
		log.Println("HTTP server listening on :" + port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Bot started")
	bot.Start()
}
