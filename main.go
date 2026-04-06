package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tele "gopkg.in/tucnak/telebot.v2"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is empty")
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot started")

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

	bot.Start()
}
