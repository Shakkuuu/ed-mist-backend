package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/service"
	"github.com/Shakkuuu/ed-mist-backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

// AppHandler アプリ向けハンドラー
type AppHandler struct {
	authUsecase          *usecase.AppAuthUsecase
	stayLogUsecase       *usecase.StayLogUsecase
	attendanceUsecase    *usecase.AttendanceUsecase
	lessonService        *service.LessonService
	deviceService        *service.DeviceService
	stayService          *service.StayService
	organizationService  *service.OrganizationService
}

// NewAppHandler アプリ向けハンドラーを作成
func NewAppHandler(
	authUsecase *usecase.AppAuthUsecase,
	stayLogUsecase *usecase.StayLogUsecase,
	attendanceUsecase *usecase.AttendanceUsecase,
	lessonService *service.LessonService,
	deviceService *service.DeviceService,
	stayService *service.StayService,
	organizationService *service.OrganizationService,
) *AppHandler {
	return &AppHandler{
		authUsecase:         authUsecase,
		stayLogUsecase:      stayLogUsecase,
		attendanceUsecase:   attendanceUsecase,
		lessonService:       lessonService,
		deviceService:       deviceService,
		stayService:         stayService,
		organizationService: organizationService,
	}
}

// DeviceRegister デバイス登録（管理画面で作成済みのユーザーが前提）
// POST /device-register
//
// リクエスト:
// - mail: メールアドレス（必須）
// - device_id: デバイスID（必須）
// - org_id: 組織ID（オプション：複数組織に同じメールアドレスがある場合に指定）
//
// レスポンス:
//   - 1つだけユーザーが見つかった場合、または org_id を指定した場合:
//     { "user": {...}, "device": {...} }
//   - 複数の組織に同じメールアドレスのユーザーが見つかった場合:
//     { "message": "...", "organizations": [...], "requires_org_id": true }
func (h *AppHandler) DeviceRegister(c echo.Context) error {
	ctx := c.Request().Context()
	var request usecase.DeviceRegisterRequest

	if err := c.Bind(&request); err != nil {
		log.Printf("[DeviceRegister] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	if request.Mail == "" {
		log.Printf("[DeviceRegister] メールアドレスが空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "mailは必須です"})
	}

	if request.DeviceID == "" {
		log.Printf("[DeviceRegister] デバイスIDが空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "device_idは必須です"})
	}

	response, err := h.authUsecase.DeviceRegister(ctx, &request)
	if err != nil {
		log.Printf("[DeviceRegister] デバイス登録エラー: %v\n", err)

		// ユーザーが見つからない場合は404
		if err == usecase.ErrorUserNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "ユーザーが見つかりません"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 複数組織が見つかった場合は、ステータスコード 300 Multiple Choices を返す
	if multipleOrgs, ok := response.(*usecase.MultipleOrganizationsResponse); ok {
		return c.JSON(http.StatusMultipleChoices, multipleOrgs)
	}

	// 正常にデバイス登録できた場合
	return c.JSON(http.StatusOK, response)
}

// GetStayLogs 滞在ログ取得
// GET /logs/stays/:org_id/:room_id/:subject_id
func (h *AppHandler) GetStayLogs(c echo.Context) error {
	ctx := c.Request().Context()

	// URLパラメータの取得
	orgID := c.Param("org_id")
	roomID := c.Param("room_id")
	subjectID := c.Param("subject_id")

	if orgID == "" {
		log.Printf("[GetStayLogs] 組織IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "組織IDが指定されていません"})
	}

	if roomID == "" {
		log.Printf("[GetStayLogs] 部屋IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "部屋IDが指定されていません"})
	}

	if subjectID == "" {
		log.Printf("[GetStayLogs] 科目IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "科目IDが指定されていません"})
	}

	// クエリパラメータの取得
	userID := c.QueryParam("user_id")
	isActiveStr := c.QueryParam("is_active")
	startTimeStr := c.QueryParam("start_time")
	endTimeStr := c.QueryParam("end_time")

	// クエリパラメータの解析
	queryReq, err := usecase.ParseStayLogsQuery(userID, isActiveStr, startTimeStr, endTimeStr)
	if err != nil {
		log.Printf("[GetStayLogs] クエリパラメータの解析エラー: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// リクエストの構築
	request := &usecase.GetStayLogsRequest{
		OrgID:     orgID,
		RoomID:    roomID,
		SubjectID: subjectID,
		UserID:    queryReq.UserID,
		IsActive:  queryReq.IsActive,
		StartTime: queryReq.StartTime,
		EndTime:   queryReq.EndTime,
	}

	stays, err := h.stayLogUsecase.GetStayLogs(ctx, request)
	if err != nil {
		log.Printf("[GetStayLogs] 滞在ログ取得エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stays)
}

// DeviceActivate デバイスアクティベーション（生体認証）
// POST /app/device/activate
func (h *AppHandler) DeviceActivate(c echo.Context) error {
	ctx := c.Request().Context()

	var request struct {
		DeviceID string `json:"device_id"`
		UserID   string `json:"user_id"`
	}

	if err := c.Bind(&request); err != nil {
		log.Printf("[DeviceActivate] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	if request.DeviceID == "" || request.UserID == "" {
		log.Printf("[DeviceActivate] device_idとuser_idは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "device_idとuser_idは必須です"})
	}

	// デバイスをアクティベーション
	device, err := h.deviceService.ActivateWithAuthentication(ctx, request.DeviceID)
	if err != nil {
		log.Printf("[DeviceActivate] デバイスアクティベーションエラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "デバイスのアクティベーションに失敗しました"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"device":  device,
		"message": "認証が完了しました",
	})
}

// RunDailyBatch 日次バッチを手動実行（デバッグ用）
// POST /app/debug/daily-batch
func (h *AppHandler) RunDailyBatch(c echo.Context) error {
	ctx := c.Request().Context()

	log.Printf("[RunDailyBatch] 日次バッチを手動実行します")

	// 全組織のデバイスを非アクティブ化
	organizations, err := h.organizationService.GetAll(ctx)
	if err != nil {
		log.Printf("[RunDailyBatch] 組織一覧取得エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "組織一覧の取得に失敗しました"})
	}

	totalDeactivated := 0
	var results []map[string]interface{}

	for _, org := range organizations {
		// 組織の全デバイスを非アクティブ化
		err := h.deviceService.DeactivateAllForOrg(ctx, org.ID)
		if err != nil {
			log.Printf("[RunDailyBatch] 組織(%s)のデバイス非アクティブ化エラー: %v", org.Name, err)
			results = append(results, map[string]interface{}{
				"organization": org.Name,
				"status":       "error",
				"error":        err.Error(),
			})
			continue
		}

		// 非アクティブ化されたデバイス数を取得
		devices, err := h.deviceService.GetAll(ctx)
		if err != nil {
			log.Printf("[RunDailyBatch] デバイス一覧取得エラー: %v", err)
			continue
		}

		orgDeviceCount := 0
		for _, device := range devices {
			if device.User.OrgID == org.ID {
				orgDeviceCount++
			}
		}

		results = append(results, map[string]interface{}{
			"organization":     org.Name,
			"status":          "success",
			"devices_count":   orgDeviceCount,
		})

		totalDeactivated += orgDeviceCount
		log.Printf("[RunDailyBatch] 組織(%s): %d台のデバイスを非アクティブ化", org.Name, orgDeviceCount)
	}

	log.Printf("[RunDailyBatch] 日次バッチ完了: %d台のデバイスを非アクティブ化", totalDeactivated)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":            true,
		"total_deactivated":  totalDeactivated,
		"organizations":      results,
		"message":            "日次バッチが完了しました",
	})
}

// GetLessonsToday 今日の時間割取得
// GET /app/lessons/today
func (h *AppHandler) GetLessonsToday(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.QueryParam("user_id")
	dateStr := c.QueryParam("date")

	if userID == "" {
		log.Printf("[GetLessonsToday] user_idは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_idは必須です"})
	}

	// 日付のパース（省略時は今日）
	var date time.Time
	var err error
	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("[GetLessonsToday] 日付の解析エラー: %v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "日付の形式が不正です（YYYY-MM-DD）"})
		}
	}

	// 時間割を取得
	lessons, err := h.lessonService.GetByUserAndDate(ctx, userID, date)
	if err != nil {
		log.Printf("[GetLessonsToday] 時間割取得エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "時間割の取得に失敗しました"})
	}

	// リレーションを整形したレスポンス
	formattedLessons := make([]map[string]interface{}, len(lessons))
	for i, lesson := range lessons {
		formattedLessons[i] = map[string]interface{}{
			"id":          lesson.ID,
			"subject_id":  lesson.SubjectID,
			"room_id":     lesson.RoomID,
			"org_id":      lesson.OrgID,
			"day_of_week": lesson.DayOfWeek,
			"start_time":  lesson.StartTime,
			"end_time":    lesson.EndTime,
			"period":      lesson.Period,
			"subject": map[string]interface{}{
				"id":   lesson.Subject.ID,
				"name": lesson.Subject.Name,
				"year": lesson.Subject.Year,
			},
			"room": map[string]interface{}{
				"id":           lesson.Room.ID,
				"org_room_id":  lesson.Room.OrgRoomID,
				"name":         lesson.Room.Name,
				"mist_zone_id": lesson.Room.MistZoneID,
			},
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"date":          date.Format("2006-01-02"),
		"day_of_week":   int(date.Weekday()),
		"lessons":       formattedLessons,
		"total_lessons": len(lessons),
	})
}

// GetAttendanceToday 今日の出席状況取得
// GET /app/attendance/today
func (h *AppHandler) GetAttendanceToday(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.QueryParam("user_id")
	dateStr := c.QueryParam("date")

	if userID == "" {
		log.Printf("[GetAttendanceToday] user_idは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_idは必須です"})
	}

	// 日付のパース（省略時は今日）
	var date time.Time
	var err error
	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("[GetAttendanceToday] 日付の解析エラー: %v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "日付の形式が不正です（YYYY-MM-DD）"})
		}
	}

	// 出席状況を取得
	records, summary, err := h.attendanceUsecase.GetTodayAttendance(ctx, userID, date)
	if err != nil {
		log.Printf("[GetAttendanceToday] 出席状況取得エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "出席状況の取得に失敗しました"})
	}

	// リレーションを整形したレスポンス
	formattedRecords := make([]map[string]interface{}, len(records))
	for i, record := range records {
		lessonData := map[string]interface{}{
			"id":          record.Lesson.ID,
			"subject_id":  record.Lesson.SubjectID,
			"room_id":     record.Lesson.RoomID,
			"day_of_week": record.Lesson.DayOfWeek,
			"start_time":  record.Lesson.StartTime,
			"end_time":    record.Lesson.EndTime,
			"period":      record.Lesson.Period,
			"subject": map[string]interface{}{
				"id":   record.Lesson.Subject.ID,
				"name": record.Lesson.Subject.Name,
				"year": record.Lesson.Subject.Year,
			},
			"room": map[string]interface{}{
				"id":           record.Lesson.Room.ID,
				"org_room_id":  record.Lesson.Room.OrgRoomID,
				"name":         record.Lesson.Room.Name,
				"mist_zone_id": record.Lesson.Room.MistZoneID,
			},
		}

		formattedRecords[i] = map[string]interface{}{
			"lesson":            lessonData,
			"attendance_status": record.AttendanceStatus,
			"late_minutes":      record.LateMinutes,
			"on_time":           record.OnTime,
			"entry_time":        record.EntryTime,
			"exit_time":         record.ExitTime,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"date":    date.Format("2006-01-02"),
		"records": formattedRecords,
		"summary": summary,
	})
}

// CreateManualStay 手動入室
// POST /app/stays/manual
func (h *AppHandler) CreateManualStay(c echo.Context) error {
	ctx := c.Request().Context()

	var request struct {
		UserID    string `json:"user_id"`
		RoomID    string `json:"room_id"`
		SubjectID string `json:"subject_id"`
	}

	if err := c.Bind(&request); err != nil {
		log.Printf("[CreateManualStay] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	if request.UserID == "" || request.RoomID == "" || request.SubjectID == "" {
		log.Printf("[CreateManualStay] user_id、room_id、subject_idは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id、room_id、subject_idは必須です"})
	}

	// 前のアクティブな滞在を終了
	previousStay, _ := h.stayService.GetActiveByUserID(ctx, request.UserID)
	if previousStay != nil {
		previousStay.IsActive = false
		now := time.Now()
		previousStay.LeavedAt = &now
		if err := h.stayService.Update(ctx, previousStay, previousStay.SubjectID, previousStay.Description); err != nil {
			log.Printf("[CreateManualStay] 前の滞在終了エラー: %v\n", err)
		}
	}

	// 現在時刻と部屋から授業を検索（自動紐付け）
	now := time.Now()
	var lessonIDPtr *string
	lesson, err := h.lessonService.GetByRoomAndTime(ctx, request.RoomID, now)
	if err == nil && lesson != nil {
		lessonIDPtr = &lesson.ID
		log.Printf("[CreateManualStay] 授業を自動検出: Lesson=%s, Subject=%s", lesson.ID, lesson.SubjectID)
	}

	// 新しい滞在を作成
	stay := &model.Stay{
		UserID:    request.UserID,
		RoomID:    request.RoomID,
		SubjectID: request.SubjectID,
		LessonID:  lessonIDPtr,
		Source:    "manual",
		IsActive:  true,
		CreatedAt: now,
	}

	if err := h.stayService.CreateWithLesson(ctx, stay); err != nil {
		log.Printf("[CreateManualStay] 滞在作成エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "入室記録に失敗しました"})
	}

	log.Printf("[CreateManualStay] 作成成功: StayID=%d, UserID=%s", stay.ID, stay.UserID)

	// Preloadを含めて再取得
	createdStay, err := h.stayService.GetByID(ctx, stay.ID)
	if err != nil {
		log.Printf("[CreateManualStay] 滞在再取得エラー: StayID=%d, Error=%v\n", stay.ID, err)
		createdStay = stay // エラー時は元のstayを使用
	} else {
		log.Printf("[CreateManualStay] 再取得成功: StayID=%d", createdStay.ID)
	}

	// リレーションを整形したレスポンス
	stayData := map[string]interface{}{
		"id":         createdStay.ID,
		"user_id":    createdStay.UserID,
		"room_id":    createdStay.RoomID,
		"subject_id": createdStay.SubjectID,
		"lesson_id":  createdStay.LessonID,
		"source":     createdStay.Source,
		"is_active":  createdStay.IsActive,
		"created_at": createdStay.CreatedAt,
		"leaved_at":  createdStay.LeavedAt,
		"room": map[string]interface{}{
			"id":           createdStay.Room.ID,
			"org_room_id":  createdStay.Room.OrgRoomID,
			"name":         createdStay.Room.Name,
			"mist_zone_id": createdStay.Room.MistZoneID,
		},
		"subject": map[string]interface{}{
			"id":   createdStay.Subject.ID,
			"name": createdStay.Subject.Name,
			"year": createdStay.Subject.Year,
		},
		"user": map[string]interface{}{
			"id":   createdStay.User.ID,
			"mail": createdStay.User.Mail,
		},
	}

	// Lessonがある場合のみ追加
	if createdStay.Lesson != nil {
		stayData["lesson"] = map[string]interface{}{
			"id":          createdStay.Lesson.ID,
			"start_time":  createdStay.Lesson.StartTime,
			"end_time":    createdStay.Lesson.EndTime,
			"period":      createdStay.Lesson.Period,
			"day_of_week": createdStay.Lesson.DayOfWeek,
			"subject": map[string]interface{}{
				"id":   createdStay.Lesson.Subject.ID,
				"name": createdStay.Lesson.Subject.Name,
			},
			"room": map[string]interface{}{
				"id":   createdStay.Lesson.Room.ID,
				"name": createdStay.Lesson.Room.Name,
			},
		}
	}

	response := map[string]interface{}{
		"stay":    stayData,
		"message": "手動入室を記録しました",
	}

	if previousStay != nil {
		response["previous_stay_closed"] = map[string]interface{}{
			"id":         previousStay.ID,
			"room_id":    previousStay.RoomID,
			"subject_id": previousStay.SubjectID,
			"leaved_at":  previousStay.LeavedAt,
		}
		response["message"] = "前の滞在を終了し、新しい入室を記録しました"
	}

	return c.JSON(http.StatusOK, response)
}

// LeaveStay 手動退室
// PUT /app/stays/:stay_id/leave
func (h *AppHandler) LeaveStay(c echo.Context) error {
	ctx := c.Request().Context()

	stayIDStr := c.Param("stay_id")
	if stayIDStr == "" {
		log.Printf("[LeaveStay] stay_idは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "stay_idは必須です"})
	}

	// stayIDを整数に変換
	var stayID int
	if _, err := fmt.Sscanf(stayIDStr, "%d", &stayID); err != nil {
		log.Printf("[LeaveStay] stay_idの解析エラー: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "stay_idの形式が不正です"})
	}

	var request struct {
		UserID string `json:"user_id"`
	}

	if err := c.Bind(&request); err != nil {
		log.Printf("[LeaveStay] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	if request.UserID == "" {
		log.Printf("[LeaveStay] user_idは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_idは必須です"})
	}

	// 滞在を取得
	stay, err := h.stayService.GetByID(ctx, stayID)
	if err != nil {
		log.Printf("[LeaveStay] 滞在取得エラー: %v\n", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "滞在ログが見つかりません"})
	}

	// ユーザーIDを確認
	if stay.UserID != request.UserID {
		log.Printf("[LeaveStay] ユーザーID不一致: %s != %s\n", stay.UserID, request.UserID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "この滞在ログにアクセスする権限がありません"})
	}

	// 退室処理
	stay.IsActive = false
	now := time.Now()
	stay.LeavedAt = &now
	if err := h.stayService.Update(ctx, stay, stay.SubjectID, stay.Description); err != nil {
		log.Printf("[LeaveStay] 退室処理エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "退室記録に失敗しました"})
	}

	// リレーションを整形したレスポンス
	stayData := map[string]interface{}{
		"id":         stay.ID,
		"user_id":    stay.UserID,
		"room_id":    stay.RoomID,
		"subject_id": stay.SubjectID,
		"lesson_id":  stay.LessonID,
		"source":     stay.Source,
		"is_active":  stay.IsActive,
		"created_at": stay.CreatedAt,
		"leaved_at":  stay.LeavedAt,
		"room": map[string]interface{}{
			"id":   stay.Room.ID,
			"name": stay.Room.Name,
		},
		"subject": map[string]interface{}{
			"id":   stay.Subject.ID,
			"name": stay.Subject.Name,
		},
		"user": map[string]interface{}{
			"id":   stay.User.ID,
			"mail": stay.User.Mail,
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stay":    stayData,
		"message": "退室を記録しました",
	})
}

// GetActiveStay アクティブな滞在確認
// GET /app/stays/active
func (h *AppHandler) GetActiveStay(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.QueryParam("user_id")
	if userID == "" {
		log.Printf("[GetActiveStay] user_idは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_idは必須です"})
	}

	// アクティブな滞在を取得
	stay, err := h.stayService.GetActiveByUserID(ctx, userID)
	if err != nil {
		// 見つからない場合はnullを返す
		return c.JSON(http.StatusOK, map[string]interface{}{
			"active_stay": nil,
		})
	}

	// 滞在時間を計算
	durationMinutes := int(time.Since(stay.CreatedAt).Minutes())

	// リレーションを整形
	activeStayData := map[string]interface{}{
		"id":               stay.ID,
		"user_id":          stay.UserID,
		"source":           stay.Source,
		"created_at":       stay.CreatedAt,
		"duration_minutes": durationMinutes,
		"room": map[string]interface{}{
			"id":           stay.Room.ID,
			"org_room_id":  stay.Room.OrgRoomID,
			"name":         stay.Room.Name,
			"mist_zone_id": stay.Room.MistZoneID,
		},
		"subject": map[string]interface{}{
			"id":   stay.Subject.ID,
			"name": stay.Subject.Name,
			"year": stay.Subject.Year,
		},
	}

	// Lessonがある場合のみ追加
	if stay.Lesson != nil {
		activeStayData["lesson"] = map[string]interface{}{
			"id":          stay.Lesson.ID,
			"start_time":  stay.Lesson.StartTime,
			"end_time":    stay.Lesson.EndTime,
			"period":      stay.Lesson.Period,
			"day_of_week": stay.Lesson.DayOfWeek,
		}
	} else {
		activeStayData["lesson"] = nil
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"active_stay": activeStayData,
	})
}

// GetUserStays ユーザーの滞在ログ取得
// GET /app/stays/:user_id
func (h *AppHandler) GetUserStays(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Param("user_id")
	if userID == "" {
		log.Printf("[GetUserStays] user_idは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_idは必須です"})
	}

	// ユーザーの滞在ログを取得
	stays, err := h.stayService.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("[GetUserStays] 滞在ログ取得エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "滞在ログの取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, stays)
}
