package main

import (
	"log"

	"github.com/Shakkuuu/ed-mist-backend/internal/app"
	"github.com/Shakkuuu/ed-mist-backend/internal/config"
	"github.com/Shakkuuu/ed-mist-backend/internal/db"

	"github.com/Shakkuuu/ed-mist-backend/pkg/mistapi"
)

func main() {
	log.SetPrefix("[APP] ")
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("設定読み込みエラー: %v", err)
	}

	// データベース接続
	dbConn, err := db.NewConnection(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}

	// Mist APIクライアントの初期化
	var mistClient *mistapi.Client
	if cfg.MistAPIToken != "" && cfg.MistSiteID != "" {
		mistClient = mistapi.NewClient(cfg.MistBaseURL, cfg.MistAPIToken, cfg.MistSiteID)
		log.Printf("Mist APIクライアントが初期化されました (SiteID: %s)\n", cfg.MistSiteID)
	} else {
		log.Println("Mist API設定が不完全なため、Mist APIクライアントは無効化されています")
	}

	app.Run(cfg, dbConn, mistClient)
}
