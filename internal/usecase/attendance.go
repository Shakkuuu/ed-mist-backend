package usecase

import (
	"context"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/service"
)

// AttendanceStatus 出席ステータス
type AttendanceStatus string

const (
	AttendanceOnTime   AttendanceStatus = "on_time"   // 定刻
	AttendanceLate     AttendanceStatus = "late"      // 遅刻
	AttendanceVeryLate AttendanceStatus = "very_late" // 大幅遅刻
	AttendanceAbsent   AttendanceStatus = "absent"    // 欠席
	AttendanceUnknown  AttendanceStatus = "unknown"   // 不明
)

// AttendanceConfig 出席管理設定
type AttendanceConfig struct {
	LateThresholdMinutes int  `json:"late_threshold_minutes"` // 遅刻許容時間（分）デフォルト: 10
	EarlyEntryMinutes    int  `json:"early_entry_minutes"`    // 授業前何分から入室可能か デフォルト: 10
	AutoCheckoutEnabled  bool `json:"auto_checkout_enabled"`  // 自動退出を有効にするか デフォルト: true
}

// DefaultAttendanceConfig デフォルトの出席管理設定
func DefaultAttendanceConfig() AttendanceConfig {
	return AttendanceConfig{
		LateThresholdMinutes: 10,
		EarlyEntryMinutes:    10,
		AutoCheckoutEnabled:  true,
	}
}

// StayWithAttendance 出席情報付き滞在ログ
type StayWithAttendance struct {
	model.Stay
	AttendanceStatus AttendanceStatus `json:"attendance_status"`
	LateMinutes      int              `json:"late_minutes"`
	OnTime           bool             `json:"on_time"`
}

// AttendanceRecord 出席記録（API レスポンス用）
type AttendanceRecord struct {
	Lesson           *model.Lesson    `json:"lesson"`
	AttendanceStatus AttendanceStatus `json:"attendance_status"`
	LateMinutes      int              `json:"late_minutes"`
	OnTime           bool             `json:"on_time"`
	EntryTime        *time.Time       `json:"entry_time,omitempty"`
	ExitTime         *time.Time       `json:"exit_time,omitempty"`
}

// AttendanceSummary 出席サマリー
type AttendanceSummary struct {
	TotalLessons   int     `json:"total_lessons"`
	OnTime         int     `json:"on_time"`
	Late           int     `json:"late"`
	Absent         int     `json:"absent"`
	AttendanceRate float64 `json:"attendance_rate"`
}

// AttendanceUsecase 出席判定ユースケース
type AttendanceUsecase struct {
	lessonService *service.LessonService
	stayService   *service.StayService
	userService   *service.UserService
}

// NewAttendanceUsecase 出席判定ユースケースを作成
func NewAttendanceUsecase(
	lessonService *service.LessonService,
	stayService *service.StayService,
	userService *service.UserService,
) *AttendanceUsecase {
	return &AttendanceUsecase{
		lessonService: lessonService,
		stayService:   stayService,
		userService:   userService,
	}
}

// CalculateAttendanceStatus 出席ステータスを計算
func CalculateAttendanceStatus(stay model.Stay, lesson *model.Lesson, config AttendanceConfig) AttendanceStatus {
	// Lessonがstayに紐付いている場合はそれを使用、なければ引数のlessonを使用
	targetLesson := stay.Lesson
	if targetLesson == nil {
		targetLesson = lesson
	}

	if targetLesson == nil {
		return AttendanceUnknown // Lessonがない
	}

	diff := stay.CreatedAt.Sub(targetLesson.StartTime)

	if diff <= 0 {
		return AttendanceOnTime // 定刻または早め
	}

	lateMinutes := int(diff.Minutes())

	if lateMinutes <= config.LateThresholdMinutes {
		return AttendanceLate // 許容範囲内の遅刻
	}

	return AttendanceVeryLate // 大幅遅刻
}

// CalculateLateMinutes 遅刻時間を計算
func CalculateLateMinutes(stay model.Stay, lesson *model.Lesson) int {
	// Lessonがstayに紐付いている場合はそれを使用、なければ引数のlessonを使用
	targetLesson := stay.Lesson
	if targetLesson == nil {
		targetLesson = lesson
	}

	if targetLesson == nil {
		return 0
	}

	diff := stay.CreatedAt.Sub(targetLesson.StartTime)
	if diff <= 0 {
		return 0
	}

	return int(diff.Minutes())
}

// EnrichStayWithAttendance 滞在ログに出席情報を付加
func EnrichStayWithAttendance(stay model.Stay, lesson *model.Lesson, config AttendanceConfig) StayWithAttendance {
	lateMinutes := CalculateLateMinutes(stay, lesson)
	status := CalculateAttendanceStatus(stay, lesson, config)

	return StayWithAttendance{
		Stay:             stay,
		AttendanceStatus: status,
		LateMinutes:      lateMinutes,
		OnTime:           lateMinutes == 0,
	}
}

// GetTodayAttendance 今日の出席状況を取得
func (u *AttendanceUsecase) GetTodayAttendance(ctx context.Context, userID string, date time.Time) ([]AttendanceRecord, *AttendanceSummary, error) {
	// ユーザー情報を取得
	user, err := u.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// 今日の授業一覧を取得
	lessons, err := u.lessonService.GetByUserAndDate(ctx, userID, date)
	if err != nil {
		return nil, nil, err
	}

	// デフォルト設定を使用（将来的には組織の設定を取得）
	config := DefaultAttendanceConfig()

	records := []AttendanceRecord{}
	summary := &AttendanceSummary{
		TotalLessons: len(lessons),
	}

	for _, lesson := range lessons {
		// この授業の滞在ログを検索
		// LessonIDで検索 + 手動入室も含める（同じ時間帯・同じ部屋）
		stays, err := u.stayService.GetByLessonID(ctx, lesson.ID)
		if err != nil {
			return nil, nil, err
		}

		// このユーザーの滞在ログを探す（LessonID一致）
		var userStay *model.Stay
		for i := range stays {
			if stays[i].UserID == user.ID {
				userStay = &stays[i]
				break
			}
		}

		// LessonIDで見つからなかった場合、時間帯と部屋で検索（手動入室対応）
		if userStay == nil {
			allUserStays, err := u.stayService.GetByUserID(ctx, user.ID)
			if err == nil {
				for i := range allUserStays {
					stay := &allUserStays[i]
					// 同じ部屋で、授業時間内に作成された滞在を探す
					if stay.RoomID == lesson.RoomID &&
						!stay.CreatedAt.Before(lesson.StartTime.Add(-10*time.Minute)) &&
						!stay.CreatedAt.After(lesson.EndTime.Add(30*time.Minute)) {
						userStay = stay
						break
					}
				}
			}
		}

		if userStay == nil {
			// 欠席
			records = append(records, AttendanceRecord{
				Lesson:           &lesson,
				AttendanceStatus: AttendanceAbsent,
				LateMinutes:      0,
				OnTime:           false,
				EntryTime:        nil,
				ExitTime:         nil,
			})
			summary.Absent++
		} else {
			// 出席
			status := CalculateAttendanceStatus(*userStay, &lesson, config)
			lateMinutes := CalculateLateMinutes(*userStay, &lesson)

			var exitTime *time.Time
			if userStay.LeavedAt != nil {
				exitTime = userStay.LeavedAt
			}

			records = append(records, AttendanceRecord{
				Lesson:           &lesson,
				AttendanceStatus: status,
				LateMinutes:      lateMinutes,
				OnTime:           lateMinutes == 0,
				EntryTime:        &userStay.CreatedAt,
				ExitTime:         exitTime,
			})

			switch status {
			case AttendanceOnTime:
				summary.OnTime++
			case AttendanceLate, AttendanceVeryLate:
				summary.Late++
			}
		}
	}

	// 出席率を計算
	if summary.TotalLessons > 0 {
		attendedLessons := summary.OnTime + summary.Late
		summary.AttendanceRate = float64(attendedLessons) / float64(summary.TotalLessons) * 100
	} else {
		summary.AttendanceRate = 0
	}

	return records, summary, nil
}
