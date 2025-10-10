package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/service"
	"github.com/Shakkuuu/ed-mist-backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

// AdminHandler 管理向けハンドラー
type AdminHandler struct {
	organizationUsecase *usecase.OrganizationUsecase
	userUsecase         *usecase.UserUsecase
	roomUsecase         *usecase.RoomUsecase
	stayLogUsecase      *usecase.StayLogUsecase
	subjectService      *service.SubjectService
	lessonService       *service.LessonService
}

// NewAdminHandler 管理向けハンドラーを作成
func NewAdminHandler(
	organizationUsecase *usecase.OrganizationUsecase,
	userUsecase *usecase.UserUsecase,
	roomUsecase *usecase.RoomUsecase,
	stayLogUsecase *usecase.StayLogUsecase,
	subjectService *service.SubjectService,
	lessonService *service.LessonService,
) *AdminHandler {
	return &AdminHandler{
		organizationUsecase: organizationUsecase,
		userUsecase:         userUsecase,
		roomUsecase:         roomUsecase,
		stayLogUsecase:      stayLogUsecase,
		subjectService:      subjectService,
		lessonService:       lessonService,
	}
}

// CreateOrganization 組織登録
// POST /organizations
func (h *AdminHandler) CreateOrganization(c echo.Context) error {
	ctx := c.Request().Context()
	var request usecase.CreateOrganizationRequest

	if err := c.Bind(&request); err != nil {
		log.Printf("[CreateOrganization] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	if request.Mail == "" {
		log.Printf("[CreateOrganization] メールアドレスが空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "mailは必須です"})
	}

	if request.Name == "" {
		log.Printf("[CreateOrganization] 組織名が空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_nameは必須です"})
	}

	organization, err := h.organizationUsecase.CreateOrganization(ctx, &request)
	if err != nil {
		log.Printf("[CreateOrganization] 組織作成エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, organization)
}

// GetOrganization 組織情報取得
// GET /organizations/:org_id
func (h *AdminHandler) GetOrganization(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")

	if orgID == "" {
		log.Printf("[GetOrganization] 組織IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "組織IDが指定されていません"})
	}

	organization, err := h.organizationUsecase.GetOrganization(ctx, orgID)
	if err != nil {
		log.Printf("[GetOrganization] 組織取得エラー: %v, orgID: %s\n", err, orgID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, organization)
}

// DeleteOrganization 組織削除
// DELETE /organizations/:org_id
func (h *AdminHandler) DeleteOrganization(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")

	if orgID == "" {
		log.Printf("[DeleteOrganization] 組織IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "組織IDが指定されていません"})
	}

	if err := h.organizationUsecase.DeleteOrganization(ctx, orgID); err != nil {
		log.Printf("[DeleteOrganization] 組織削除エラー: %v, orgID: %s\n", err, orgID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "組織が削除されました"})
}

// CreateUser ユーザー（追跡対象者）登録
// POST /users
func (h *AdminHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	var request usecase.CreateUserRequest

	if err := c.Bind(&request); err != nil {
		log.Printf("[CreateUser] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	if request.OrgID == "" {
		log.Printf("[CreateUser] 組織IDが空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_idは必須です"})
	}

	if request.UserMail == "" {
		log.Printf("[CreateUser] ユーザーメールが空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_mailは必須です"})
	}

	user, err := h.userUsecase.CreateUser(ctx, &request)
	if err != nil {
		log.Printf("[CreateUser] ユーザー作成エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// GetUsers ユーザー一覧
// GET /users/:org_id
func (h *AdminHandler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")

	if orgID == "" {
		log.Printf("[GetUsers] 組織IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "組織IDが指定されていません"})
	}

	users, err := h.userUsecase.GetUsersByOrgID(ctx, orgID)
	if err != nil {
		log.Printf("[GetUsers] ユーザー一覧取得エラー: %v, orgID: %s\n", err, orgID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

// GetUser ユーザー情報取得
// GET /users/:org_id/:user_id
func (h *AdminHandler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")
	userID := c.Param("user_id")

	if orgID == "" {
		log.Printf("[GetUser] 組織IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "組織IDが指定されていません"})
	}

	if userID == "" {
		log.Printf("[GetUser] ユーザーIDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ユーザーIDが指定されていません"})
	}

	user, err := h.userUsecase.GetUser(ctx, orgID, userID)
	if err != nil {
		log.Printf("[GetUser] ユーザー取得エラー: %v, orgID: %s, userID: %s\n", err, orgID, userID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser ユーザー削除
// DELETE /users/:org_id/:user_id
func (h *AdminHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")
	userID := c.Param("user_id")

	if orgID == "" {
		log.Printf("[DeleteUser] 組織IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "組織IDが指定されていません"})
	}

	if userID == "" {
		log.Printf("[DeleteUser] ユーザーIDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ユーザーIDが指定されていません"})
	}

	if err := h.userUsecase.DeleteUser(ctx, orgID, userID); err != nil {
		log.Printf("[DeleteUser] ユーザー削除エラー: %v, orgID: %s, userID: %s\n", err, orgID, userID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ユーザーが削除されました"})
}

// CreateRoom 部屋情報登録
// POST /rooms
func (h *AdminHandler) CreateRoom(c echo.Context) error {
	ctx := c.Request().Context()
	var request usecase.CreateRoomRequest

	if err := c.Bind(&request); err != nil {
		log.Printf("[CreateRoom] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	if request.OrgID == "" {
		log.Printf("[CreateRoom] 組織IDが空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_idは必須です"})
	}

	if request.OrgRoomID == "" {
		log.Printf("[CreateRoom] 組織部屋IDが空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_room_idは必須です"})
	}

	if request.RoomName == "" {
		log.Printf("[CreateRoom] 部屋名が空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "room_nameは必須です"})
	}

	room, err := h.roomUsecase.CreateRoom(ctx, &request)
	if err != nil {
		log.Printf("[CreateRoom] 部屋作成エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, room)
}

// GetRooms 部屋一覧取得
// GET /rooms/:org_id
func (h *AdminHandler) GetRooms(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")

	if orgID == "" {
		log.Printf("[GetRooms] 組織IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "組織IDが指定されていません"})
	}

	rooms, err := h.roomUsecase.GetRoomsByOrgID(ctx, orgID)
	if err != nil {
		log.Printf("[GetRooms] 部屋一覧取得エラー: %v, orgID: %s\n", err, orgID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, rooms)
}

// UpdateRoom 部屋情報更新
// PUT /rooms/:org_id/:room_id
func (h *AdminHandler) UpdateRoom(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")
	roomID := c.Param("room_id")
	var request usecase.UpdateRoomRequest

	if orgID == "" {
		log.Printf("[UpdateRoom] 組織IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "組織IDが指定されていません"})
	}

	if roomID == "" {
		log.Printf("[UpdateRoom] 部屋IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "部屋IDが指定されていません"})
	}

	if err := c.Bind(&request); err != nil {
		log.Printf("[UpdateRoom] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	room, err := h.roomUsecase.UpdateRoom(ctx, orgID, roomID, &request)
	if err != nil {
		log.Printf("[UpdateRoom] 部屋更新エラー: %v, orgID: %s, roomID: %s\n", err, orgID, roomID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, room)
}

// DeleteRoom 部屋情報削除
// DELETE /rooms/:org_id/:room_id
func (h *AdminHandler) DeleteRoom(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")
	roomID := c.Param("room_id")

	if orgID == "" {
		log.Printf("[DeleteRoom] 組織IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "組織IDが指定されていません"})
	}

	if roomID == "" {
		log.Printf("[DeleteRoom] 部屋IDが指定されていません\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "部屋IDが指定されていません"})
	}

	if err := h.roomUsecase.DeleteRoom(ctx, orgID, roomID); err != nil {
		log.Printf("[DeleteRoom] 部屋削除エラー: %v, orgID: %s, roomID: %s\n", err, orgID, roomID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "部屋が削除されました"})
}

// GetStayLogs 滞在ログ取得（管理向け）
// GET /logs/stays/:org_id/:room_id/:subject_id
func (h *AdminHandler) GetStayLogs(c echo.Context) error {
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

// CreateSubject 教科作成
// POST /api/v1/subjects
func (h *AdminHandler) CreateSubject(c echo.Context) error {
	ctx := c.Request().Context()

	var request struct {
		OrgID string `json:"org_id"`
		Name  string `json:"name"`
		Year  int    `json:"year"`
	}

	if err := c.Bind(&request); err != nil {
		log.Printf("[CreateSubject] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	if request.OrgID == "" || request.Name == "" {
		log.Printf("[CreateSubject] org_id, nameは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_id, nameは必須です"})
	}

	if request.Year == 0 {
		request.Year = time.Now().Year()
	}

	subject, err := h.subjectService.Create(ctx, request.Name, request.Year, request.OrgID)
	if err != nil {
		log.Printf("[CreateSubject] 教科作成エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "教科の作成に失敗しました"})
	}

	// リレーションを除外したレスポンス
	response := map[string]interface{}{
		"id":         subject.ID,
		"name":       subject.Name,
		"year":       subject.Year,
		"org_id":     subject.OrgID,
		"created_at": subject.CreatedAt,
		"updated_at": subject.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, response)
}

// GetSubjects 教科一覧取得
// GET /api/v1/subjects/:org_id
func (h *AdminHandler) GetSubjects(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")

	if orgID == "" {
		log.Printf("[GetSubjects] org_idは必須です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_idは必須です"})
	}

	yearStr := c.QueryParam("year")
	var err error

	var subjects interface{}
	if yearStr != "" {
		year, parseErr := strconv.Atoi(yearStr)
		if parseErr != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "yearの形式が不正です"})
		}
		subjects, err = h.subjectService.GetByOrgIDAndYear(ctx, orgID, year)
	} else {
		subjects, err = h.subjectService.GetByOrgID(ctx, orgID)
	}

	if err != nil {
		log.Printf("[GetSubjects] 教科一覧取得エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "教科一覧の取得に失敗しました"})
	}

	// リレーションなしのクリーンなレスポンス（既にRepository層で除外済み）
	return c.JSON(http.StatusOK, subjects)
}

// CreateLesson 授業作成
// POST /api/v1/lessons
func (h *AdminHandler) CreateLesson(c echo.Context) error {
	ctx := c.Request().Context()

	var request struct {
		OrgID      string `json:"org_id"`
		SubjectID  string `json:"subject_id"`
		RoomID     string `json:"room_id"`
		DayOfWeek  int    `json:"day_of_week"`
		StartTime  string `json:"start_time"`  // "09:00"
		EndTime    string `json:"end_time"`    // "10:30"
		Period     int    `json:"period"`
		DateString string `json:"date"`        // "2025-10-10" (オプション)
	}

	if err := c.Bind(&request); err != nil {
		log.Printf("[CreateLesson] リクエストの解析に失敗しました: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
	}

	if request.OrgID == "" || request.SubjectID == "" || request.RoomID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_id, subject_id, room_idは必須です"})
	}

	if request.StartTime == "" || request.EndTime == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "start_time, end_timeは必須です"})
	}

	// 日付を取得（指定されていなければ今日）
	var baseDate time.Time
	if request.DateString != "" {
		parsedDate, err := time.Parse("2006-01-02", request.DateString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "dateの形式が不正です（YYYY-MM-DD）"})
		}
		baseDate = parsedDate
	} else {
		baseDate = time.Now()
	}

	// 時刻をパース
	startTimeParsed, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "start_timeの形式が不正です（HH:MM）"})
	}

	endTimeParsed, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "end_timeの形式が不正です（HH:MM）"})
	}

	// baseDateに時刻を設定
	startTime := time.Date(
		baseDate.Year(), baseDate.Month(), baseDate.Day(),
		startTimeParsed.Hour(), startTimeParsed.Minute(), 0, 0,
		time.UTC,
	)
	endTime := time.Date(
		baseDate.Year(), baseDate.Month(), baseDate.Day(),
		endTimeParsed.Hour(), endTimeParsed.Minute(), 0, 0,
		time.UTC,
	)

	// day_of_weekが指定されていなければ、baseDateから計算
	dayOfWeek := request.DayOfWeek
	if dayOfWeek == 0 {
		dayOfWeek = int(baseDate.Weekday())
	}

	var datePtr *time.Time
	if request.DateString != "" {
		datePtr = &baseDate
	}

	lesson, err := h.lessonService.Create(
		ctx,
		request.SubjectID,
		request.RoomID,
		request.OrgID,
		dayOfWeek,
		startTime,
		endTime,
		request.Period,
		datePtr,
	)
	if err != nil {
		log.Printf("[CreateLesson] 授業作成エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "授業の作成に失敗しました"})
	}

	// リレーションを除外したレスポンス
	response := map[string]interface{}{
		"id":          lesson.ID,
		"subject_id":  lesson.SubjectID,
		"room_id":     lesson.RoomID,
		"org_id":      lesson.OrgID,
		"day_of_week": lesson.DayOfWeek,
		"start_time":  lesson.StartTime,
		"end_time":    lesson.EndTime,
		"period":      lesson.Period,
		"date":        lesson.Date,
		"created_at":  lesson.CreatedAt,
		"updated_at":  lesson.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, response)
}

// GetLessons 授業一覧取得
// GET /api/v1/lessons/:org_id
func (h *AdminHandler) GetLessons(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("org_id")

	if orgID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_idは必須です"})
	}

	lessons, err := h.lessonService.GetByOrgID(ctx, orgID)
	if err != nil {
		log.Printf("[GetLessons] 授業一覧取得エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "授業一覧の取得に失敗しました"})
	}

	// リレーションを整形したレスポンス
	response := make([]map[string]interface{}, len(lessons))
	for i, lesson := range lessons {
		response[i] = map[string]interface{}{
			"id":          lesson.ID,
			"subject_id":  lesson.SubjectID,
			"room_id":     lesson.RoomID,
			"org_id":      lesson.OrgID,
			"day_of_week": lesson.DayOfWeek,
			"start_time":  lesson.StartTime,
			"end_time":    lesson.EndTime,
			"period":      lesson.Period,
			"date":        lesson.Date,
			"created_at":  lesson.CreatedAt,
			"updated_at":  lesson.UpdatedAt,
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

	return c.JSON(http.StatusOK, response)
}

// DeleteLesson 授業削除
// DELETE /api/v1/lessons/:lesson_id
func (h *AdminHandler) DeleteLesson(c echo.Context) error {
	ctx := c.Request().Context()
	lessonID := c.Param("lesson_id")

	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "lesson_idは必須です"})
	}

	if err := h.lessonService.Delete(ctx, lessonID); err != nil {
		log.Printf("[DeleteLesson] 授業削除エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "授業の削除に失敗しました"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "授業が削除されました"})
}
