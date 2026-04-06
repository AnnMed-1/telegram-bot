package main

import (
	"fmt"
	"log"
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

	bot, err := tele.NewBot(tele.Settings{
		Token: token,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot started")

	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Отправь ссылку")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		url := c.Text()
		_ = c.Notify(tele.Typing)

		tmpDir := os.TempDir()
		pattern := filepath.Join(tmpDir, "video.%(ext)s")

		cmd := exec.Command("yt-dlp", "-o", pattern, url)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return c.Send(fmt.Sprintf("Ошибка запуска команды:\n%s", string(out)))
		}

		matches, _ := filepath.Glob(filepath.Join(tmpDir, "video.*"))
		if len(matches) == 0 {
			return c.Send("Файл скачался, но не найден")
		}

		videoPath := matches[0]
		video := &tele.Video{File: tele.FromDisk(videoPath)}

		if err := c.Send(video); err != nil {
			return c.Send("Не удалось отправить видео: " + err.Error())
		}

		_ = os.Remove(videoPath)
		return nil
	})

	bot.Start()
}
