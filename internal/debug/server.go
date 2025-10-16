package debug

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Shakkuuu/ed-mist-backend/internal/debug/handler"
	"github.com/Shakkuuu/ed-mist-backend/internal/debug/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DebugServer デバッグ用サーバー
type DebugServer struct {
	db     *gorm.DB
	server *echo.Echo
}

// NewDebugServer デバッグサーバーを作成
func NewDebugServer() *DebugServer {
	return &DebugServer{}
}

// Start デバッグサーバーを起動
func (s *DebugServer) Start() error {
	// データベース接続
	if err := s.initDatabase(); err != nil {
		return fmt.Errorf("データベース初期化エラー: %w", err)
	}

	// Echoサーバーの初期化
	s.server = echo.New()
	s.setupMiddleware()
	s.setupRoutes()

	// サーバー起動
	port := s.getPort()
	log.Printf("デバッグサーバーを起動中: :%d", port)

	return s.server.Start(":" + strconv.Itoa(port))
}

// initDatabase データベースを初期化
func (s *DebugServer) initDatabase() error {
	// 環境変数からデータベース接続情報を取得
	// host := getEnv("DB_HOST", "localhost")
	host := "localhost"
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER_NAME", "postgres")
	password := getEnv("DB_USER_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "mist_ed")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// GORMロガーの設定
	newLogger := logger.New(
		log.New(os.Stdout, "[GORM] ", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)

	// データベース接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("データベース接続エラー: %w", err)
	}

	// 接続テスト
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("データベース接続テストエラー: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("データベース接続テストエラー: %w", err)
	}

	s.db = db
	log.Println("データベース接続が確立されました")
	return nil
}

// setupMiddleware ミドルウェアを設定
func (s *DebugServer) setupMiddleware() {
	s.server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[DEBUG] ${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}\n",
	}))
	s.server.Use(middleware.Recover())
	s.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
}

// setupRoutes ルートを設定
func (s *DebugServer) setupRoutes() {
	// リポジトリの初期化
	orgRepo := repository.NewOrganizationRepository(s.db)
	userRepo := repository.NewUserRepository(s.db)
	roomRepo := repository.NewRoomRepository(s.db)
	deviceRepo := repository.NewDeviceRepository(s.db)
	subjectRepo := repository.NewSubjectRepository(s.db)
	lessonRepo := repository.NewLessonRepository(s.db)

	// ハンドラーの初期化
	debugHandler := handler.NewDebugHandler(orgRepo, userRepo, roomRepo, deviceRepo, subjectRepo, lessonRepo, s.db)

	// デバッグ用ルート
	debug := s.server.Group("/debug")
	debugHandler.RegisterRoutes(debug)

	// ヘルスチェック
	s.server.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok", "service": "debug-server"})
	})
}

// getPort ポート番号を取得
func (s *DebugServer) getPort() int {
	portStr := getEnv("DEBUG_SERVER_PORT", "8081")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("警告: 無効なポート番号 '%s', デフォルトの8081を使用します", portStr)
		return 8081
	}
	return port
}

// getEnv 環境変数を取得（デフォルト値付き）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
