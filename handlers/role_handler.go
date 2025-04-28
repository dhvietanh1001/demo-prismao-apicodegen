// handlers/role_handler.go
package handlers

import (
	"context"
	"demo-prismao-apicodegen"
	"demo-prismao-apicodegen/prisma/db"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RoleHandler xử lý các request cho Role endpoints
type RoleHandler struct {
	Client *db.PrismaClient
}

// GetRoles trả về danh sách tất cả roles
func (h *RoleHandler) GetRoles(ctx echo.Context) error {
	prismCtx := context.Background()
	roles, err := h.Client.Role.FindMany().Exec(prismCtx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, roles)
}

// GetRolesId trả về một role dựa theo ID
func (h *RoleHandler) GetRolesId(ctx echo.Context, id int) error {
	prismCtx := context.Background()
	role, err := h.Client.Role.FindUnique(
		db.Role.ID.Equals(id),
	).Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if role == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
	}

	return ctx.JSON(http.StatusOK, role)
}

// PostRoles tạo một role mới
func (h *RoleHandler) PostRoles(ctx echo.Context) error {
	var roleCreate api.RoleCreate
	if err := ctx.Bind(&roleCreate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	prismCtx := context.Background()
	role, err := h.Client.Role.CreateOne(
		db.Role.Name.Set(roleCreate.Name),
		db.Role.Description.SetIfPresent(roleCreate.Description),
		db.Role.AssignedTo.SetIfPresent(roleCreate.AssignedTo),
	).Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, role)
}

// DeleteRolesId xóa một role
func (h *RoleHandler) DeleteRolesId(ctx echo.Context, id int) error {
	prismCtx := context.Background()
	_, err := h.Client.Role.FindUnique(
		db.Role.ID.Equals(id),
	).Delete().Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// PutRolesId cập nhật thông tin role
func (h *RoleHandler) PutRolesId(ctx echo.Context, id int) error {
	var roleUpdate api.RoleUpdate
	if err := ctx.Bind(&roleUpdate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	prismCtx := context.Background()
	params := []db.RoleSetParam{}

	if roleUpdate.Name != nil {
		params = append(params, db.Role.Name.Set(*roleUpdate.Name))
	}

	if roleUpdate.Description != nil {
		params = append(params, db.Role.Description.Set(*roleUpdate.Description))
	}

	if roleUpdate.AssignedTo != nil {
		params = append(params, db.Role.AssignedTo.Set(*roleUpdate.AssignedTo))
	}

	role, err := h.Client.Role.FindUnique(
		db.Role.ID.Equals(id),
	).Update(
		params...,
	).Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, role)
}
