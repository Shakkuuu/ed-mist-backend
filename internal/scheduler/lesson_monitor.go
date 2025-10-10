package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

// LessonMonitor 授業監視ワーカー
type LessonMonitor struct {
	lesson        model.Lesson
	scheduler     *LessonScheduler
	recordedUsers map[string]bool // すでに記録したユーザー
	stopChan      chan struct{}
}

// NewLessonMonitor 授業監視ワーカーを作成
func NewLessonMonitor(lesson model.Lesson, scheduler *LessonScheduler) *LessonMonitor {
	return &LessonMonitor{
		lesson:        lesson,
		scheduler:     scheduler,
		recordedUsers: make(map[string]bool),
		stopChan:      make(chan struct{}),
	}
}

// Start 監視を開始
func (m *LessonMonitor) Start() {
	defer m.cleanup()

	// 監視期間
	monitorStart := m.lesson.StartTime.Add(-5 * time.Minute)
	monitorEnd := m.lesson.EndTime.Add(10 * time.Minute)

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	log.Printf("[LessonMonitor] 監視開始: Lesson=%s, 期間=%s〜%s",
		m.lesson.ID,
		monitorStart.Format("15:04"),
		monitorEnd.Format("15:04"))

	// 即座に1回チェック
	m.checkZone()

	for {
		select {
		case <-m.stopChan:
			log.Printf("[LessonMonitor] 停止: Lesson=%s", m.lesson.ID)
			return
		case <-ticker.C:
			now := time.Now()

			// 監視終了チェック
			if now.After(monitorEnd) {
				log.Printf("[LessonMonitor] 監視終了: Lesson=%s", m.lesson.ID)
				m.finishLesson()
				return
			}

			// 監視開始前なら待機
			if now.Before(monitorStart) {
				continue
			}

			// Zone監視実行
			m.checkZone()
		}
	}
}

// Stop 監視を停止
func (m *LessonMonitor) Stop() {
	close(m.stopChan)
}

// checkZone Zone内のデバイスをチェック
func (m *LessonMonitor) checkZone() {
	ctx := context.Background()

	// Room情報を取得
	room, err := m.scheduler.roomService.GetByID(ctx, m.lesson.RoomID)
	if err != nil {
		log.Printf("[LessonMonitor] Room取得エラー: %v", err)
		return
	}

	if room.MistZoneID == "" {
		log.Printf("[LessonMonitor] Room(%s)にMistZoneIDが設定されていません", room.ID)
		return
	}

	// Mist APIでZone内のデバイス一覧を取得
	sdkClients, wirelessClients, err := m.scheduler.mistClient.GetZoneClients(
		m.scheduler.mistClient.SiteID,
		room.MistZoneID,
	)
	if err != nil {
		log.Printf("[LessonMonitor] Zone監視エラー: %v", err)
		return
	}

	allDevices := append(sdkClients, wirelessClients...)
	log.Printf("[LessonMonitor] 検知デバイス数: %d (Lesson=%s, Zone=%s)",
		len(allDevices), m.lesson.ID, room.MistZoneID)

	// 各デバイスについて処理
	for _, deviceID := range allDevices {
		m.processDevice(deviceID)
	}
}

// processDevice デバイスを処理して出席記録
func (m *LessonMonitor) processDevice(deviceID string) {
	ctx := context.Background()

	// デバイスIDからユーザーIDを取得
	device, err := m.scheduler.deviceService.GetByDeviceID(ctx, deviceID)
	if err != nil {
		// デバイスが登録されていない
		return
	}

	userID := device.UserID

	// すでに記録済みかチェック
	if m.recordedUsers[userID] {
		return
	}

	// 本日認証済みかチェック
	if !device.IsAuthenticatedToday() {
		log.Printf("[LessonMonitor] 未認証デバイス: User=%s, Device=%s", userID, deviceID)
		return
	}

	// 滞在ログを作成
	now := time.Now()
	lessonID := m.lesson.ID
	stay := &model.Stay{
		UserID:    userID,
		RoomID:    m.lesson.RoomID,
		SubjectID: m.lesson.SubjectID,
		LessonID:  &lessonID,
		Source:    "auto",
		IsActive:  true,
		CreatedAt: now,
	}

	err = m.scheduler.stayService.CreateWithLesson(ctx, stay)
	if err != nil {
		log.Printf("[LessonMonitor] 滞在ログ作成エラー: %v", err)
		return
	}

	// 記録済みとしてマーク
	m.recordedUsers[userID] = true

	// 遅刻判定
	lateMinutes := 0
	if now.After(m.lesson.StartTime) {
		lateMinutes = int(now.Sub(m.lesson.StartTime).Minutes())
	}

	if lateMinutes > 0 {
		log.Printf("[LessonMonitor] 出席記録（遅刻）: User=%s, Lesson=%s, Time=%s, Late=%dmin",
			userID, m.lesson.ID, now.Format("15:04:05"), lateMinutes)
	} else {
		log.Printf("[LessonMonitor] 出席記録（定刻）: User=%s, Lesson=%s, Time=%s",
			userID, m.lesson.ID, now.Format("15:04:05"))
	}
}

// finishLesson 授業を終了
func (m *LessonMonitor) finishLesson() {
	ctx := context.Background()

	log.Printf("[LessonMonitor] 授業終了処理開始: Lesson=%s, 出席者数=%d",
		m.lesson.ID, len(m.recordedUsers))

	// すべての滞在ログを終了
	for userID := range m.recordedUsers {
		stay, err := m.scheduler.stayService.GetActiveByUserAndLesson(ctx, userID, m.lesson.ID)
		if err != nil || stay == nil {
			continue
		}

		stay.IsActive = false
		now := time.Now()
		stay.LeavedAt = &now

		err = m.scheduler.stayService.Update(ctx, stay, stay.SubjectID, stay.Description)
		if err != nil {
			log.Printf("[LessonMonitor] 退出処理エラー: User=%s, %v", userID, err)
		} else {
			log.Printf("[LessonMonitor] 自動退出: User=%s, Lesson=%s", userID, m.lesson.ID)
		}
	}

	log.Printf("[LessonMonitor] 授業終了処理完了: Lesson=%s", m.lesson.ID)
}

// cleanup クリーンアップ
func (m *LessonMonitor) cleanup() {
	m.scheduler.removeMonitor(m.lesson.ID)
	log.Printf("[LessonMonitor] クリーンアップ完了: Lesson=%s", m.lesson.ID)
}
