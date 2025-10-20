package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/config"
	"github.com/Shakkuuu/ed-mist-backend/internal/db"
	"github.com/Shakkuuu/ed-mist-backend/internal/handler"
	"github.com/Shakkuuu/ed-mist-backend/internal/repository"
	"github.com/Shakkuuu/ed-mist-backend/internal/scheduler"
	"github.com/Shakkuuu/ed-mist-backend/internal/service"
	"github.com/Shakkuuu/ed-mist-backend/internal/usecase"
	"github.com/Shakkuuu/ed-mist-backend/pkg/mistapi"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(cfg *config.Config, dbConn *db.Connection, mistClient *mistapi.Client) {
	log.SetPrefix("[APP] ")

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(dbConn.DB)
	deviceRepo := repository.NewDeviceRepository(dbConn.DB)
	organizationRepo := repository.NewOrganizationRepository(dbConn.DB)
	roomRepo := repository.NewRoomRepository(dbConn.DB)
	stayRepo := repository.NewStayRepository(dbConn.DB)
	subjectRepo := repository.NewSubjectRepository(dbConn.DB)
	lessonRepo := repository.NewLessonRepository(dbConn.DB)

	// serviceの初期化
	userService := service.NewUserService(userRepo)
	deviceService := service.NewDeviceService(deviceRepo)
	organizationService := service.NewOrganizationService(organizationRepo)
	roomService := service.NewRoomService(roomRepo)
	stayService := service.NewStayService(stayRepo)
	subjectService := service.NewSubjectService(subjectRepo)
	lessonService := service.NewLessonService(lessonRepo)
	zoneService := service.NewZoneService(mistClient)
	mapService := service.NewMapService(mistClient)

	// usecaseの初期化
	organizationUsecase := usecase.NewOrganizationUsecase(organizationService)
	userUsecase := usecase.NewUserUsecase(userService, organizationService)
	roomUsecase := usecase.NewRoomUsecase(roomService, organizationService)
	appAuthUsecase := usecase.NewAppAuthUsecase(userService, deviceService, organizationService)
	stayLogUsecase := usecase.NewStayLogUsecase(stayService, userService, roomService, subjectService, organizationService)
	attendanceUsecase := usecase.NewAttendanceUsecase(lessonService, stayService, userService)

	// APIハンドラーの初期化
	appHandler := handler.NewAppHandler(appAuthUsecase, stayLogUsecase, attendanceUsecase, lessonService, deviceService, stayService, organizationService)
	adminHandler := handler.NewAdminHandler(organizationUsecase, userUsecase, roomUsecase, stayLogUsecase, subjectService, lessonService)

	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[ECHO] ${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human} bytes_in=${bytes_in} bytes_out=${bytes_out}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// サービス情報をログに出力
	log.Printf("初期化されたサービス:")
	log.Printf("- UserService: %v", userService != nil)
	log.Printf("- DeviceService: %v", deviceService != nil)
	log.Printf("- OrganizationService: %v", organizationService != nil)
	log.Printf("- RoomService: %v", roomService != nil)
	log.Printf("- StayService: %v", stayService != nil)
	log.Printf("- SubjectService: %v", subjectService != nil)
	log.Printf("- ZoneService: %v", zoneService != nil)
	log.Printf("- MapService: %v", mapService != nil)

	// 授業スケジューラーの初期化と起動
	lessonScheduler := scheduler.NewLessonScheduler(
		lessonService,
		roomService,
		deviceService,
		stayService,
		mistClient,
	)
	go lessonScheduler.Start()
	log.Println("授業スケジューラーを起動しました")

	// 日次バッチスケジューラーの初期化と起動
	dailyBatchScheduler := scheduler.NewDailyBatchScheduler(
		deviceService,
		organizationService,
	)
	go dailyBatchScheduler.Start()
	log.Println("日次バッチスケジューラーを起動しました")

	// API
	apiV1 := e.Group("/api/v1")
	{
		// アプリ向けエンドポイント
		app := e.Group("/app")
		{
			// デバイス登録（管理画面で作成済みのユーザーが前提）
			app.POST("/device-register", appHandler.DeviceRegister)

			// デバイスアクティベーション（生体認証）
			app.POST("/device/activate", appHandler.DeviceActivate)

			// 時間割取得
			app.GET("/lessons/today", appHandler.GetLessonsToday)

			// 出席状況取得
			app.GET("/attendance/today", appHandler.GetAttendanceToday)

			// 手動入室
			app.POST("/stays/manual", appHandler.CreateManualStay)

			// 手動退室
			app.PUT("/stays/:stay_id/leave", appHandler.LeaveStay)

			// アクティブな滞在確認
			app.GET("/stays/active", appHandler.GetActiveStay)

			// 滞在ログ取得（ユーザー向け）
			app.GET("/stays/:user_id", appHandler.GetUserStays)

			// デバッグ用エンドポイント
			app.POST("/debug/daily-batch", appHandler.RunDailyBatch)
		}

		// 管理向けエンドポイント
		// 組織関連
		organizations := apiV1.Group("/organizations")
		{
			organizations.POST("", adminHandler.CreateOrganization)
			organizations.GET("", adminHandler.GetOrganizations)
			organizations.GET("/:org_id", adminHandler.GetOrganization)
			organizations.DELETE("/:org_id", adminHandler.DeleteOrganization)
		}

		// ユーザー関連
		users := apiV1.Group("/users")
		{
			users.POST("", adminHandler.CreateUser)
			users.GET("/:org_id", adminHandler.GetUsers)
			users.GET("/:org_id/:user_id", adminHandler.GetUser)
			users.DELETE("/:org_id/:user_id", adminHandler.DeleteUser)
		}

		// 部屋関連
		rooms := apiV1.Group("/rooms")
		{
			rooms.POST("", adminHandler.CreateRoom)
			rooms.GET("/:org_id", adminHandler.GetRooms)
			rooms.PUT("/:org_id/:room_id", adminHandler.UpdateRoom)
			rooms.DELETE("/:org_id/:room_id", adminHandler.DeleteRoom)
		}

		// 滞在ログ取得（管理向け）
		logs := apiV1.Group("/logs")
		{
			logs.GET("/stays/:org_id/:room_id/:subject_id", adminHandler.GetStayLogs)
		}

		// 教科関連
		subjects := apiV1.Group("/subjects")
		{
			subjects.POST("", adminHandler.CreateSubject)
			subjects.GET("/:org_id", adminHandler.GetSubjects)
			subjects.DELETE("/:subject_id", adminHandler.DeleteSubject)
		}

		// 授業関連
		lessons := apiV1.Group("/lessons")
		{
			lessons.POST("", adminHandler.CreateLesson)
			lessons.GET("/:org_id", adminHandler.GetLessons)
			lessons.DELETE("/:lesson_id", adminHandler.DeleteLesson)
		}
	}

	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// シグナルチャンネルの作成
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// サーバー起動
	addr := ":" + strconv.Itoa(cfg.ServerPort)
	log.Printf("サーバーを起動中: %s\n", addr)

	// goroutineでサーバーを起動
	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバー起動エラー: %v", err)
		}
	}()

	// シグナル待機
	<-quit
	log.Println("Graceful Shutdownを開始します...")

	// スケジューラーを停止
	log.Println("授業スケジューラーを停止しています...")
	lessonScheduler.Stop()

	log.Println("日次バッチスケジューラーを停止しています...")
	dailyBatchScheduler.Stop()

	// タイムアウト付きのcontextでシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("サーバーシャットダウンエラー: %v", err)
	}

	log.Println("サーバーが正常にシャットダウンされました")
}
