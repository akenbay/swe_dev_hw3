package handler

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"university/internal/middleware"
	"university/internal/model"
	"university/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// Register registers HTTP routes on the provided Echo instance.
func (h *Handler) Register(e *echo.Echo) {
	// Public auth routes
	e.POST("/api/auth/register", h.Register_User)
	e.POST("/api/auth/login", h.Login)

	// Protected routes
	e.GET("/api/users/me", h.GetCurrentUser, middleware.AuthMiddleware(h.service))

	// Public student/schedule routes
	e.GET("/student/:id", h.GetStudentByID)
	e.GET("/students", h.GetAllStudents)
	e.GET("/all_class_schedule", h.GetAllSchedules)
	e.GET("/schedule/group/:id", h.GetGroupSchedule)
}

func (h *Handler) GetStudentByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	student, err := h.service.GetStudentByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "student not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, student)
}

func (h *Handler) GetAllStudents(c echo.Context) error {
	students, err := h.service.GetAllStudents()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, students)
}

func (h *Handler) GetAllSchedules(c echo.Context) error {
	schedules, err := h.service.GetAllSchedules()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, schedules)
}

func (h *Handler) GetGroupSchedule(c echo.Context) error {
	groupID := c.Param("id")
	if groupID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "group id is required"})
	}

	schedules, err := h.service.GetGroupSchedule(groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, schedules)
}

func (h *Handler) CreateAttendanceRecord(c echo.Context) error {
	var record model.AttendanceRecord
	if err := c.Bind(&record); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err := h.service.CreateAttendanceRecord(&record)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "attendance record created successfully"})
}

// Register_User handles user registration
func (h *Handler) Register_User(c echo.Context) error {
	var req model.AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	user, err := h.service.Register(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Login handles user login and returns JWT token
func (h *Handler) Login(c echo.Context) error {
	var req model.AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	response, err := h.service.Login(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// GetCurrentUser returns current user info (protected endpoint)
func (h *Handler) GetCurrentUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user context"})
	}

	user, err := h.service.GetCurrentUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetAttendanceRecordsByStudentID(c echo.Context) error {
	studentID := c.Param("id")
	if studentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	records, err := h.service.GetAttendanceRecordsByStudentID(studentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, records)
}

func (h *Handler) GetAttendanceRecordsBySubjectID(c echo.Context) error {
	subjectID := c.Param("id")
	if subjectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	records, err := h.service.GetAttendanceRecordsBySubjectID(subjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, records)
}
