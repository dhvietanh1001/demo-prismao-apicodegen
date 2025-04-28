// handlers/user_handler.go
package handlers

import (
	"context"
	"demo-prismao-apicodegen"
	"demo-prismao-apicodegen/prisma/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// UserHandler xử lý các request cho User endpoints
type UserHandler struct {
	Client *db.PrismaClient
}

// GetUsers trả về danh sách tất cả users
func (h *UserHandler) GetUsers(ctx echo.Context) error {
	prismCtx := context.Background()
	users, err := h.Client.User.FindMany().Exec(prismCtx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, users)
}

// GetUsersId trả về một user dựa theo ID
func (h *UserHandler) GetUsersId(ctx echo.Context, id int) error {
	prismCtx := context.Background()
	user, err := h.Client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if user == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return ctx.JSON(http.StatusOK, user)
}

// PostUsers tạo một user mới
func (h *UserHandler) PostUsers(ctx echo.Context) error {
	var userCreate api.UserCreate
	if err := ctx.Bind(&userCreate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	prismCtx := context.Background()
	user, err := h.Client.User.CreateOne(
		db.User.Name.Set(userCreate.Name),
		db.User.Email.Set(string(userCreate.Email)),
		db.User.PasswordHash.Set(string(hashedPassword)),
		db.User.SsoUserid.SetIfPresent(userCreate.SsoUserid),
	).Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, user)
}

// DeleteUsersId xóa một user
func (h *UserHandler) DeleteUsersId(ctx echo.Context, id int) error {
	prismCtx := context.Background()
	_, err := h.Client.User.FindUnique(
		db.User.ID.Equals(id),
	).Delete().Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// PutUsersId cập nhật thông tin user
func (h *UserHandler) PutUsersId(ctx echo.Context, id int) error {
	var userUpdate api.UserUpdate
	if err := ctx.Bind(&userUpdate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	prismCtx := context.Background()
	params := []db.UserSetParam{}

	if userUpdate.Name != nil {
		params = append(params, db.User.Name.Set(*userUpdate.Name))
	}

	if userUpdate.Email != nil {
		params = append(params, db.User.Email.Set(string(*userUpdate.Email)))
	}

	if userUpdate.Password != nil {
		// Trong thực tế, bạn nên hash password trước khi lưu
		params = append(params, db.User.PasswordHash.Set(*userUpdate.Password))
	}

	if userUpdate.SsoUserid != nil {
		params = append(params, db.User.SsoUserid.Set(*userUpdate.SsoUserid))
	}

	user, err := h.Client.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		params...,
	).Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}
