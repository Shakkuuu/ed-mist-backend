package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/service"
	"github.com/Shakkuuu/ed-mist-backend/pkg/mistapi"
)

// LessonScheduler 授業スケジューラー
type LessonScheduler struct {
	lessonService *service.LessonService
	roomService   *service.RoomService
	deviceService *service.DeviceService
	stayService   *service.StayService
	mistClient    *mistapi.Client

	activeMonitors sync.Map // map[lessonID]*LessonMonitor
	stopChan       chan struct{}
}

// NewLessonScheduler 授業スケジューラーを作成
func NewLessonScheduler(
	lessonService *service.LessonService,
	roomService *service.RoomService,
	deviceService *service.DeviceService,
	stayService *service.StayService,
	mistClient *mistapi.Client,
) *LessonScheduler {
	return &LessonScheduler{
		lessonService: lessonService,
		roomService:   roomService,
		deviceService: deviceService,
		stayService:   stayService,
		mistClient:    mistClient,
		stopChan:      make(chan struct{}),
	}
}

// Start スケジューラーを開始
func (s *LessonScheduler) Start() {
	log.Println("[LessonScheduler] 開始")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			log.Println("[LessonScheduler] 停止")
			return
		case <-ticker.C:
			s.checkAndStartMonitors()
		}
	}
}

// Stop スケジューラーを停止
func (s *LessonScheduler) Stop() {
	close(s.stopChan)
}

// checkAndStartMonitors 監視対象の授業をチェックして監視を開始
func (s *LessonScheduler) checkAndStartMonitors() {
	ctx := context.Background()
	now := time.Now()

	// 現在時刻で監視対象の授業を取得
	lessons, err := s.lessonService.GetMonitoringLessons(ctx, now)
	if err != nil {
		log.Printf("[LessonScheduler] 監視対象授業取得エラー: %v", err)
		return
	}

	log.Printf("[LessonScheduler] 監視対象授業数: %d", len(lessons))

	// 監視対象の授業IDのセットを作成
	monitoringLessonIDs := make(map[string]bool)
	for _, lesson := range lessons {
		monitoringLessonIDs[lesson.ID] = true
	}

	// 既存の監視プロセスをチェック
	s.activeMonitors.Range(func(key, value interface{}) bool {
		lessonID := key.(string)
		monitor := value.(*LessonMonitor)

		// 監視対象から外れた場合は停止
		if !monitoringLessonIDs[lessonID] {
			log.Printf("[LessonScheduler] 監視停止: Lesson=%s", lessonID)
			monitor.Stop()
			s.activeMonitors.Delete(lessonID)
		}
		return true
	})

	for _, lesson := range lessons {
		// すでに監視中かチェック
		if _, exists := s.activeMonitors.Load(lesson.ID); exists {
			continue
		}

		// 新しい授業の監視を開始
		log.Printf("[LessonScheduler] 監視開始: Lesson=%s, Subject=%s, Room=%s, Time=%s-%s",
			lesson.ID, lesson.SubjectID, lesson.RoomID,
			lesson.StartTime.Format("15:04"), lesson.EndTime.Format("15:04"))

		monitor := NewLessonMonitor(lesson, s)
		s.activeMonitors.Store(lesson.ID, monitor)

		go monitor.Start()
	}
}

// removeMonitor 監視を削除
func (s *LessonScheduler) removeMonitor(lessonID string) {
	s.activeMonitors.Delete(lessonID)
}
