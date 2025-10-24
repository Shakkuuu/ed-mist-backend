package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/debug/repository"
	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// DebugHandler デバッグ用ハンドラー
type DebugHandler struct {
	orgRepo     *repository.OrganizationRepository
	userRepo    *repository.UserRepository
	roomRepo    *repository.RoomRepository
	deviceRepo  *repository.DeviceRepository
	subjectRepo *repository.SubjectRepository
	lessonRepo  *repository.LessonRepository
	db          *gorm.DB
}

// NewDebugHandler デバッグハンドラーを作成
func NewDebugHandler(
	orgRepo *repository.OrganizationRepository,
	userRepo *repository.UserRepository,
	roomRepo *repository.RoomRepository,
	deviceRepo *repository.DeviceRepository,
	subjectRepo *repository.SubjectRepository,
	lessonRepo *repository.LessonRepository,
	db *gorm.DB,
) *DebugHandler {
	return &DebugHandler{
		orgRepo:     orgRepo,
		userRepo:    userRepo,
		roomRepo:    roomRepo,
		deviceRepo:  deviceRepo,
		subjectRepo: subjectRepo,
		lessonRepo:  lessonRepo,
		db:          db,
	}
}

// RegisterRoutes デバッグ用ルートを登録
func (h *DebugHandler) RegisterRoutes(e *echo.Group) {
	// GET endpoints
	e.GET("/organizations", h.GetOrganizations)
	e.GET("/users", h.GetUsers)
	e.GET("/rooms", h.GetRooms)
	e.GET("/devices", h.GetDevices)
	e.GET("/subjects", h.GetSubjects)
	e.GET("/lessons", h.GetLessons)

	// POST endpoints
	e.POST("/organizations", h.CreateOrganization)
	e.POST("/users", h.CreateUser)
	e.POST("/rooms", h.CreateRoom)
	e.POST("/devices", h.CreateDevice)
	e.POST("/subjects", h.CreateSubject)
	e.POST("/lessons", h.CreateLesson)

	// DELETE endpoints
	e.DELETE("/organizations", h.DeleteOrganizations)
	e.DELETE("/users", h.DeleteUsers)
	e.DELETE("/rooms", h.DeleteRooms)
	e.DELETE("/devices", h.DeleteDevices)
	e.DELETE("/subjects", h.DeleteSubjects)
	e.DELETE("/lessons", h.DeleteLessons)

	// Special endpoints
	e.POST("/seed", h.CreateSeedData)
	e.DELETE("/reset", h.ResetDatabase)
}

// GetOrganizations 組織一覧取得
func (h *DebugHandler) GetOrganizations(c echo.Context) error {
	orgs, err := h.orgRepo.FindAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, orgs)
}

// GetUsers ユーザー一覧取得
func (h *DebugHandler) GetUsers(c echo.Context) error {
	users, err := h.userRepo.FindAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

// GetRooms 部屋一覧取得
func (h *DebugHandler) GetRooms(c echo.Context) error {
	rooms, err := h.roomRepo.FindAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rooms)
}

// GetDevices デバイス一覧取得
func (h *DebugHandler) GetDevices(c echo.Context) error {
	devices, err := h.deviceRepo.FindAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, devices)
}

// GetSubjects 科目一覧取得
func (h *DebugHandler) GetSubjects(c echo.Context) error {
	subjects, err := h.subjectRepo.FindAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, subjects)
}

// GetLessons 授業一覧取得
func (h *DebugHandler) GetLessons(c echo.Context) error {
	lessons, err := h.lessonRepo.FindAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, lessons)
}

// CreateOrganization 組織作成
func (h *DebugHandler) CreateOrganization(c echo.Context) error {
	var req struct {
		Name string `json:"name"`
		Mail string `json:"mail"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	org := &model.Organization{
		ID:   uuid.New().String(),
		Mail: req.Mail,
		Name: req.Name,
	}

	if err := h.orgRepo.Create(context.Background(), org); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, org)
}

// CreateUser ユーザー作成
func (h *DebugHandler) CreateUser(c echo.Context) error {
	var req struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		OrganizationID string `json:"organization_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	user := &model.User{
		ID:    uuid.New().String(),
		OrgID: req.OrganizationID,
		Mail:  req.Email,
	}

	if err := h.userRepo.Create(context.Background(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// CreateRoom 部屋作成
func (h *DebugHandler) CreateRoom(c echo.Context) error {
	var req struct {
		Name           string `json:"name"`
		OrgRoomID      string `json:"org_room_id"`
		Caption        string `json:"caption"`
		MistZoneID     string `json:"mist_zone_id"`
		MapID          string `json:"map_id"`
		OrganizationID string `json:"organization_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	room := &model.Room{
		ID:         uuid.New().String(),
		OrgID:      req.OrganizationID,
		OrgRoomID:  req.OrgRoomID,
		Name:       req.Name,
		Caption:    req.Caption,
		MistZoneID: req.MistZoneID,
		MapID:      req.MapID,
	}

	if err := h.roomRepo.Create(context.Background(), room); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, room)
}

// CreateDevice デバイス作成
func (h *DebugHandler) CreateDevice(c echo.Context) error {
	var req struct {
		UserID   string `json:"user_id"`
		DeviceID string `json:"device_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	device := &model.Device{
		ID:                uuid.New().String(),
		UserID:            req.UserID,
		DeviceID:          req.DeviceID,
		IsActive:          true,
		LastAuthenticated: time.Now(),
	}

	if err := h.deviceRepo.Create(context.Background(), device); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, device)
}

// CreateSubject 科目作成
func (h *DebugHandler) CreateSubject(c echo.Context) error {
	var req struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		OrganizationID string `json:"organization_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	subject := &model.Subject{
		ID:    uuid.New().String(),
		OrgID: req.OrganizationID,
		Name:  req.Name,
	}

	if err := h.subjectRepo.Create(context.Background(), subject); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, subject)
}

// CreateLesson 授業作成
func (h *DebugHandler) CreateLesson(c echo.Context) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		SubjectID   string `json:"subject_id"`
		RoomID      string `json:"room_id"`
		OrgID       string `json:"org_id"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// 必須フィールドのバリデーション
	if req.OrgID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_id is required"})
	}
	if req.SubjectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "subject_id is required"})
	}
	if req.RoomID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "room_id is required"})
	}

	startTime, err := time.Parse("2006-01-02T15:04", req.StartTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_time format"})
	}

	endTime, err := time.Parse("2006-01-02T15:04", req.EndTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_time format"})
	}

	lesson := &model.Lesson{
		ID:        uuid.New().String(),
		SubjectID: req.SubjectID,
		RoomID:    req.RoomID,
		OrgID:     req.OrgID,
		DayOfWeek: int(startTime.Weekday()),
		StartTime: startTime,
		EndTime:   endTime,
		Period:    1,
	}

	if err := h.lessonRepo.Create(context.Background(), lesson); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, lesson)
}

// DeleteOrganizations 組織削除
func (h *DebugHandler) DeleteOrganizations(c echo.Context) error {
	if err := h.orgRepo.DeleteAll(context.Background()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "全組織を削除しました"})
}

// DeleteUsers ユーザー削除
func (h *DebugHandler) DeleteUsers(c echo.Context) error {
	if err := h.userRepo.DeleteAll(context.Background()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "全ユーザーを削除しました"})
}

// DeleteRooms 部屋削除
func (h *DebugHandler) DeleteRooms(c echo.Context) error {
	if err := h.roomRepo.DeleteAll(context.Background()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "全部屋を削除しました"})
}

// DeleteDevices デバイス削除
func (h *DebugHandler) DeleteDevices(c echo.Context) error {
	if err := h.deviceRepo.DeleteAll(context.Background()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "全デバイスを削除しました"})
}

// DeleteSubjects 科目削除
func (h *DebugHandler) DeleteSubjects(c echo.Context) error {
	if err := h.subjectRepo.DeleteAll(context.Background()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "全科目を削除しました"})
}

// DeleteLessons 授業削除
func (h *DebugHandler) DeleteLessons(c echo.Context) error {
	if err := h.lessonRepo.DeleteAll(context.Background()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "全授業を削除しました"})
}

// CreateSeedData シードデータ作成
func (h *DebugHandler) CreateSeedData(c echo.Context) error {
	// 組織を作成
	org := &model.Organization{
		ID:   uuid.New().String(),
		Mail: "seed-org@example.com",
		Name: "シード組織",
	}
	if err := h.orgRepo.Create(context.Background(), org); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// ユーザーを作成
	users := []*model.User{
		{ID: uuid.New().String(), OrgID: org.ID, Mail: "student1@example.com"},
		{ID: uuid.New().String(), OrgID: org.ID, Mail: "student2@example.com"},
		{ID: uuid.New().String(), OrgID: org.ID, Mail: "teacher1@example.com"},
	}

	for _, user := range users {
		if err := h.userRepo.Create(context.Background(), user); err != nil {
			continue
		}
	}

	// 部屋を作成
	rooms := []*model.Room{
		{ID: uuid.New().String(), OrgID: org.ID, OrgRoomID: "room-101", Name: "101教室", Caption: "1階の教室"},
		{ID: uuid.New().String(), OrgID: org.ID, OrgRoomID: "room-201", Name: "201教室", Caption: "2階の教室"},
		{ID: uuid.New().String(), OrgID: org.ID, OrgRoomID: "room-301", Name: "301教室", Caption: "3階の教室"},
	}

	for _, room := range rooms {
		if err := h.roomRepo.Create(context.Background(), room); err != nil {
			continue
		}
	}

	// 科目を作成
	subjects := []*model.Subject{
		{ID: uuid.New().String(), OrgID: org.ID, Name: "数学"},
		{ID: uuid.New().String(), OrgID: org.ID, Name: "英語"},
		{ID: uuid.New().String(), OrgID: org.ID, Name: "国語"},
	}

	for _, subject := range subjects {
		if err := h.subjectRepo.Create(context.Background(), subject); err != nil {
			continue
		}
	}

	// デバイスを作成
	for i, user := range users {
		device := &model.Device{
			ID:                uuid.New().String(),
			UserID:            user.ID,
			DeviceID:          fmt.Sprintf("device-%d", i+1),
			IsActive:          true,
			LastAuthenticated: time.Now(),
		}
		if err := h.deviceRepo.Create(context.Background(), device); err != nil {
			continue
		}
	}

	// 授業を作成
	now := time.Now()
	lessons := []*model.Lesson{
		{
			ID:        uuid.New().String(),
			SubjectID: subjects[0].ID,
			RoomID:    rooms[0].ID,
			OrgID:     org.ID,
			DayOfWeek: 1, // 月曜日
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location()),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 10, 30, 0, 0, now.Location()),
			Period:    1,
		},
		{
			ID:        uuid.New().String(),
			SubjectID: subjects[1].ID,
			RoomID:    rooms[1].ID,
			OrgID:     org.ID,
			DayOfWeek: 1, // 月曜日
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, now.Location()),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 12, 30, 0, 0, now.Location()),
			Period:    2,
		},
	}

	for _, lesson := range lessons {
		if err := h.lessonRepo.Create(context.Background(), lesson); err != nil {
			continue
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "シードデータの作成が完了しました"})
}

// ResetDatabase データベースリセット
func (h *DebugHandler) ResetDatabase(c echo.Context) error {
	tables := []string{"devices", "lessons", "users", "rooms", "subjects", "organizations"}

	for _, table := range tables {
		if err := h.db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "データベースのリセットが完了しました"})
}
