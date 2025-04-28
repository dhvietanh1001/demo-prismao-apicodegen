// handlers/user_role_handler.go
package handlers

import (
	"context"
	"demo-prismao-apicodegen"
	"demo-prismao-apicodegen/prisma/db"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// UserRoleHandler xử lý các request cho UserRole endpoints
type UserRoleHandler struct {
	Client *db.PrismaClient
}

// GetUserroles trả về danh sách tất cả user roles
func (h *UserRoleHandler) GetUserroles(ctx echo.Context) error {
	prismCtx := context.Background()
	userRoles, err := h.Client.UserRole.FindMany().With(
		db.UserRole.User.Fetch(),
		db.UserRole.Role.Fetch(),
	).Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, userRoles)
}

// PostUserroles gán một role cho user
func (h *UserRoleHandler) PostUserroles(ctx echo.Context) error {
	var userRoleCreate api.UserRoleCreate
	if err := ctx.Bind(&userRoleCreate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	prismCtx := context.Background()
	userRole, err := h.Client.UserRole.CreateOne(
		db.UserRole.User.Link(
			db.User.ID.Equals(userRoleCreate.UserId),
		),
		db.UserRole.Role.Link(
			db.Role.ID.Equals(userRoleCreate.RoleId),
		),
	).Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, userRole)
}

// DeleteUserrolesId xóa một user role e
func (h *UserRoleHandler) DeleteUserrolesId(ctx echo.Context, id int) error {
	prismCtx := context.Background()

	_, err := h.Client.UserRole.FindFirst(
		db.UserRole.ID.Equals(id),
	).Exec(prismCtx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "UserRole not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	_, err = h.Client.UserRole.FindMany(
		db.UserRole.ID.Equals(id),
	).Delete().Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.NoContent(http.StatusNoContent)
}
