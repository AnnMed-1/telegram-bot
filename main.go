package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "os/exec"
    
    "gopkg.in/telebot.v3"
)

var bot *telebot.Bot

func main() {
    var err error
    bot, err = telebot.NewBot(telebot.Settings{
        Token:  "YOUR_BOT_TOKEN",
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal(err)
        return
    }
    
    bot.Handle("/start", func(c telebot.Context) error {
        return c.Send("🤖 Skillspace Downloader готов!\n📎 Отправь Skillspace ссылку")
    })
    
    bot.Handle(telebot.OnPhoto, func(m telebot.Context) error {
        // Скачиваем по QR-коду
        photo := m.Message().Photo.File
        fileInfo, _ := photo.File()
        filePath := filepath.Join("temp", fmt.Sprintf("%d.jpg", time.Now().Unix()))
        os.MkdirAll("temp", 0755)
        
        bot.Send(m.Sender, "⏳ Считываю QR...")
        telebot.DownloadFile(bot, fileInfo.FileID, filePath)
        
        cmd := exec.Command("zbarimg", "--raw", filePath)
        urlBytes, err := cmd.Output()
        if err != nil {
            bot.Send(m.Sender, "❌ QR не найден")
            os.Remove(filePath)
            return nil
        }
        
        url := string(bytes.TrimSpace(urlBytes))
        if !strings.Contains(url, "skillspace") {
            bot.Send(m.Sender, "❌ QR не Skillspace")
            os.Remove(filePath)
            return nil
        }
        
        bot.Send(m.Sender, "🔗 Найден Skillspace! Скачиваю...")
        downloadVideo(m.Sender, url)
        os.Remove(filePath)
        return nil
    })
    
    bot.Handle(telebot.OnText, func(m telebot.Context) error {
        text := m.Text()
        if strings.Contains(text, "skillspace") {
            bot.Send(m.Sender, "🔗 Skillspace найден! Скачиваю...")
            downloadVideo(m.Sender, text)
        }
        return nil
    })
    
    log.Println("✅ Skillspace Bot started!")
    bot.Start()
}

func downloadVideo(to telebot.Recipient, url string) {
    videoFile := fmt.Sprintf("temp/video_%d.mp4", time.Now().Unix())
    os.MkdirAll("temp", 0755)
    
    cmd := exec.Command("yt-dlp", "-f", "best[height<=720]", 
                       "-o", videoFile, url)
    err := cmd.Run()
    
    if err != nil {
        bot.Send(to, "❌ Ошибка скачивания")
        return
    }
    
    fileInfo, err := os.Stat(videoFile)
    if err != nil || fileInfo.Size() == 0 {
        bot.Send(to, "❌ Файл пустой или ошибка")
        return
    }
    
    // Отправляем видео (макс 50MB)
    if fileInfo.Size() > 50*1024*1024 {
        bot.Send(to, "❌ Видео >50MB, ссылаюсь на файл")
        bot.Send(to, &telebot.Document{
            File: telebot.FromDisk(videoFile),
        })
    } else {
        bot.Send(to, &telebot.Video{
            File: telebot.FromDisk(videoFile),
        })
    }
    
    // Удаляем
    os.Remove(videoFile)
}


