package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func main() {
	// 監視するディレクトリのパスを指定
	directoryPath := "/path/to/your/image/directory"

	// fsnotifyを開始
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// ディレクトリを監視
	err = filepath.Walk(directoryPath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("監視を開始しました...")

	// 変更を検知してメールを送信
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				// 新しいファイルが作成された場合
				fileInfo, err := os.Stat(event.Name)
				if err == nil && !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), ".jpg") {
					// JPGファイルの場合のみ処理
					log.Printf("新しい画像ファイルが作成されました: %s\n", event.Name)
					sendEmail("新しい画像ファイルが作成されました", "新しい画像ファイルが作成されました: "+event.Name)
				}
			}
		case err := <-watcher.Errors:
			log.Println("エラー:", err)
		}
	}
}

func sendEmail(subject, body string) {
	// メールの設定
	e := email.NewEmail()
	e.From = "送信元メールアドレス"
	e.To = []string{"宛先メールアドレス1", "宛先メールアドレス2"} // 複数の宛先メールアドレスを指定
	e.Subject = subject
	e.Text = []byte(body)

	// SMTPSサーバーの設定
	err := e.Send("hogehoge.com:465", smtp.PlainAuth("", "送信元メールアドレス", "送信元メールアカウントのパスワード", "smtp.gmail.com"))
	if err != nil {
		log.Println("メールの送信に失敗しました:", err)
	} else {
		log.Println("メールを送信しました:", body)
	}
}