package main

import (
	"log"

	"github.com/Shakkuuu/ed-mist-backend/internal/debug"
	"github.com/joho/godotenv"
)

func main() {
	log.SetPrefix("[DEBUG-SERVER] ")

	// .envファイルを読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: .envファイルの読み込みに失敗しました: %v", err)
	}

	// デバッグサーバーを起動
	debugServer := debug.NewDebugServer()
	if err := debugServer.Start(); err != nil {
		log.Fatalf("デバッグサーバー起動エラー: %v", err)
	}
}
