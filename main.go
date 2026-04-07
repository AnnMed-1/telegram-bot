package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
    
    tele "gopkg.in/tucnak/telebot.v2"
)

func main() {
    token := os.Getenv("BOT_TOKEN")
    log.Printf("Запуск с токеном: %s...", token[:10])
    
    bot, err := tele.NewBot(tele.Settings{
        Token:  token,
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal("BOT ERROR:", err)
        return
    }

    // /start
    bot.Handle("/start", func(m *tele.Message) {
        bot.Send(m.Sender, "✅ @Kurses_skil_bot готов!\n\n"+
            "📎 Отправь ссылку на Skillspace урок\n"+
            "🎥 Скачаю видео в MP4")
    })

    // Skillspace ссылки
    bot.Handle(tele.OnText, func(m *tele.Message) {
        url := m.Text
        if strings.Contains(url, "skillspace.ru") {
            m.Notify(tele.Typing)
            
            // Скачиваем
            pattern := filepath.Join("/tmp", "skillspace.%(ext)s")
            cmd := exec.Command("yt-dlp", 
                "-f", "best[ext=mp4]/best", 
                "--no-playlist", 
                "-o", pattern, 
                url)
            
            output, err := cmd.CombinedOutput()
            if err != nil {
                log.Printf("yt-dlp error: %v", err)
                bot.Send(m.Sender, "❌ Ошибка скачивания:\n"+string(output))
                return
            }
            
            
            // Найти файл
            files, _ := filepath.Glob("/tmp/skillspace.*")
            if len(files) == 0 {
                bot.Send(m.Sender, "❌ Видео не найдено")
                return
            }
            
            videoFile := files[0]
            fileInfo, err := os.Stat(videoFile)
            if err != nil || fileInfo.Size() == 0 {
                bot.Send(m.Sender, "❌ Файл пустой или ошибка")
                return
            }
            
            // Отправляем видео (макс 50MB)
            if fileInfo.Size() > 50*1024*1024 {
                bot.Send(m.Sender, "❌ Видео >50MB, ссылаюсь на файл")
                bot.Send(m.Sender, tele.File{
                    File: tele.FromDisk(videoFile),
                })
            } else {
                bot.Send(m.Sender, &tele.Video{
                    File: tele.FromDisk(videoFile),
                })
            }
            
            // Удаляем
            os.Remove(videoFile)
        } else {
            bot.Send(m.Sender, "📎 Отправь Skillspace ссылку")
        }
    })

    log.Println("✅ Skillspace Bot started!")
    bot.Start()
}
