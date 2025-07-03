package v1

import (
	"github.com/azoma13/auth-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })

	auth := handler.Group("/auth")
	{
		newAuthRoutes(auth, services.Auth)
	}
	authMiddleware := &AuthMiddleware{services.Auth}
	v1 := handler.Group("/api/v1", authMiddleware.UserIdentity)
	{
		newAccountRoutes(v1.Group("/accounts"), services.Account)
	}
}
