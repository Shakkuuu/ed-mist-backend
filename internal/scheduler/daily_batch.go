package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/service"
)

// DailyBatchScheduler 日次バッチスケジューラー
type DailyBatchScheduler struct {
	deviceService       *service.DeviceService
	organizationService *service.OrganizationService
	stopChan            chan struct{}
}

// NewDailyBatchScheduler 日次バッチスケジューラーを作成
func NewDailyBatchScheduler(
	deviceService *service.DeviceService,
	organizationService *service.OrganizationService,
) *DailyBatchScheduler {
	return &DailyBatchScheduler{
		deviceService:       deviceService,
		organizationService: organizationService,
		stopChan:            make(chan struct{}),
	}
}

// Start 日次バッチを開始
func (d *DailyBatchScheduler) Start() {
	log.Println("[DailyBatchScheduler] 日次バッチスケジューラーを開始しました")

	// 初回実行を翌日の深夜0時に設定
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	initialDelay := nextMidnight.Sub(now)

	log.Printf("[DailyBatchScheduler] 初回実行予定: %s (%.1f時間後)",
		nextMidnight.Format("2006-01-02 15:04:05"), initialDelay.Hours())

	// 初回実行まで待機
	time.Sleep(initialDelay)

	// 毎日深夜0時に実行
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// 初回実行
	d.runDailyBatch()

	for {
		select {
		case <-d.stopChan:
			log.Println("[DailyBatchScheduler] 日次バッチスケジューラーを停止しました")
			return
		case <-ticker.C:
			d.runDailyBatch()
		}
	}
}

// Stop 日次バッチを停止
func (d *DailyBatchScheduler) Stop() {
	close(d.stopChan)
}

// runDailyBatch 日次バッチを実行
func (d *DailyBatchScheduler) runDailyBatch() {
	ctx := context.Background()
	startTime := time.Now()

	log.Printf("[DailyBatchScheduler] 日次バッチ開始: %s", startTime.Format("2006-01-02 15:04:05"))

	// 全組織のデバイスを非アクティブ化
	organizations, err := d.organizationService.GetAll(ctx)
	if err != nil {
		log.Printf("[DailyBatchScheduler] 組織一覧取得エラー: %v", err)
		return
	}

	totalDeactivated := 0
	for _, org := range organizations {
		// 組織の全デバイスを非アクティブ化
		err := d.deviceService.DeactivateAllForOrg(ctx, org.ID)
		if err != nil {
			log.Printf("[DailyBatchScheduler] 組織(%s)のデバイス非アクティブ化エラー: %v", org.Name, err)
			continue
		}

		// 非アクティブ化されたデバイス数を取得（ログ用）
		devices, err := d.deviceService.GetAll(ctx)
		if err != nil {
			log.Printf("[DailyBatchScheduler] デバイス一覧取得エラー: %v", err)
			continue
		}

		orgDeviceCount := 0
		for _, device := range devices {
			if device.User.OrgID == org.ID {
				orgDeviceCount++
			}
		}

		log.Printf("[DailyBatchScheduler] 組織(%s): %d台のデバイスを非アクティブ化", org.Name, orgDeviceCount)
		totalDeactivated += orgDeviceCount
	}

	duration := time.Since(startTime)
	log.Printf("[DailyBatchScheduler] 日次バッチ完了: %d台のデバイスを非アクティブ化 (実行時間: %.2f秒)",
		totalDeactivated, duration.Seconds())
}
