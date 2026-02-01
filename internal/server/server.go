package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"university/internal/handler"

	_ "university/docs"
)

type Server struct {
	handler *handler.Handler
}

func NewServer(handler *handler.Handler) *Server {
	return &Server{handler: handler}
}

func (s *Server) Start(addr string) error {
	e := echo.New()

	// Add global middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())

	// CORS: allow browser requests from any origin (e.g. frontend on different domain)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Register all routes (public and protected)
	s.handler.Register(e)

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e.Start(addr)
}
