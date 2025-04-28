package handlers

import (
	"context"
	"demo-prismao-apicodegen"
	"demo-prismao-apicodegen/prisma/db"
	"demo-prismao-apicodegen/ultis"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// SettingHandler xử lý các request cho Setting endpoints
type SettingHandler struct {
	Client *db.PrismaClient
}

// GetSettings trả về danh sách tất cả settings
func (h *SettingHandler) GetSettings(ctx echo.Context) error {
	prismCtx := context.Background()
	settings, err := h.Client.Setting.FindMany().Exec(prismCtx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, settings)
}

// GetSettingsUserId trả về settings cho một user
func (h *SettingHandler) GetSettingsUserId(ctx echo.Context, userId int) error {
	prismCtx := context.Background()
	setting, err := h.Client.Setting.FindUnique(
		db.Setting.UserID.Equals(userId),
	).Exec(prismCtx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if setting == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Settings not found"})
	}

	return ctx.JSON(http.StatusOK, setting)
}

// PostSettingsUserId tạo settings cho một user
func (h *SettingHandler) PostSettingsUserId(ctx echo.Context, userId int) error {
	var settingCreate api.SettingCreate
	if err := ctx.Bind(&settingCreate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	preferencesData := settingCreate.Preferences
	mergedPreferences := utils.MergePreferences(utils.DefaultPreferences, preferencesData)
	preferencesJSON, err := json.Marshal(mergedPreferences)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to encode preferences"})
	}
	prismCtx := context.Background()

	// Tạo setting mới
	setting, err := h.Client.Setting.CreateOne(
		db.Setting.Preferences.Set(preferencesJSON), // Đặt preferences là JSON đã mã hóa
		db.Setting.User.Link(
			db.User.ID.Equals(userId),
		),
	).Exec(prismCtx)

	if err != nil {
		if strings.Contains(err.Error(), "Unique constraint failed") {
			return ctx.JSON(http.StatusConflict, map[string]string{"error": "Settings already exist for this user"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, setting)
}

// PutSettingsUserId cập nhật settings cho một user
func (h *SettingHandler) PutSettingsUserId(ctx echo.Context, userId int) error {
	var settingUpdate api.SettingUpdate
	if err := ctx.Bind(&settingUpdate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	preferencesData := settingUpdate.Preferences
	mergedPreferences := utils.MergePreferences(utils.DefaultPreferences, preferencesData)
	preferencesJSON, err := json.Marshal(mergedPreferences)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to encode preferences"})
	}
	prismCtx := context.Background()
	setting, err := h.Client.Setting.FindUnique(
		db.Setting.UserID.Equals(userId),
	).Update(
		db.Setting.Preferences.Set(preferencesJSON),
	).Exec(prismCtx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Settings not found for this user"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, setting)
}
