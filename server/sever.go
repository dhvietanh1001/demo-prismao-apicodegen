// server/server.go
package server

import (
	"demo-prismao-apicodegen"
	"demo-prismao-apicodegen/handlers"
	"demo-prismao-apicodegen/prisma/db"

	"github.com/labstack/echo/v4"
)

type APIServer struct {
	userHandler     *handlers.UserHandler
	roleHandler     *handlers.RoleHandler
	userRoleHandler *handlers.UserRoleHandler
	settingHandler  *handlers.SettingHandler
}

// NewAPIServer tạo một APIServer mới với Prisma client
func NewAPIServer(client *db.PrismaClient) *APIServer {
	return &APIServer{
		userHandler:     &handlers.UserHandler{Client: client},
		roleHandler:     &handlers.RoleHandler{Client: client},
		userRoleHandler: &handlers.UserRoleHandler{Client: client},
		settingHandler:  &handlers.SettingHandler{Client: client},
	}
}

func (s *APIServer) GetUsers(ctx echo.Context) error {
	return s.userHandler.GetUsers(ctx)
}

func (s *APIServer) GetUsersId(ctx echo.Context, id int) error {
	return s.userHandler.GetUsersId(ctx, id)
}

func (s *APIServer) PostUsers(ctx echo.Context) error {
	return s.userHandler.PostUsers(ctx)
}

func (s *APIServer) PutUsersId(ctx echo.Context, id int) error {
	return s.userHandler.PutUsersId(ctx, id)
}

func (s *APIServer) DeleteUsersId(ctx echo.Context, id int) error {
	return s.userHandler.DeleteUsersId(ctx, id)
}

// Role endpoints
func (s *APIServer) GetRoles(ctx echo.Context) error {
	return s.roleHandler.GetRoles(ctx)
}

func (s *APIServer) GetRolesId(ctx echo.Context, id int) error {
	return s.roleHandler.GetRolesId(ctx, id)
}

func (s *APIServer) PostRoles(ctx echo.Context) error {
	return s.roleHandler.PostRoles(ctx)
}

func (s *APIServer) PutRolesId(ctx echo.Context, id int) error {
	return s.roleHandler.PutRolesId(ctx, id)
}

func (s *APIServer) DeleteRolesId(ctx echo.Context, id int) error {
	return s.roleHandler.DeleteRolesId(ctx, id)
}

// UserRole endpoints
func (s *APIServer) GetUserroles(ctx echo.Context) error {
	return s.userRoleHandler.GetUserroles(ctx)
}

func (s *APIServer) PostUserroles(ctx echo.Context) error {
	return s.userRoleHandler.PostUserroles(ctx)
}

func (s *APIServer) DeleteUserrolesId(ctx echo.Context, id int) error {
	return s.userRoleHandler.DeleteUserrolesId(ctx, id)
}

// Setting endpoints
func (s *APIServer) GetSettings(ctx echo.Context) error {
	return s.settingHandler.GetSettings(ctx)
}

func (s *APIServer) GetSettingsUserId(ctx echo.Context, userId int) error {
	return s.settingHandler.GetSettingsUserId(ctx, userId)
}

func (s *APIServer) PostSettingsUserId(ctx echo.Context, userId int) error {
	return s.settingHandler.PostSettingsUserId(ctx, userId)
}

func (s *APIServer) PutSettingsUserId(ctx echo.Context, userId int) error {
	return s.settingHandler.PutSettingsUserId(ctx, userId)
}

func SetupEcho(client *db.PrismaClient) *echo.Echo {
	e := echo.New()
	server := NewAPIServer(client)
	api.RegisterHandlers(e, server)

	return e
}
