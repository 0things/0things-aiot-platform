package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"0things-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func deviceError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, repository.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
		status = http.StatusNotFound
	} else if errors.Is(err, repository.ErrVersionConflict) {
		status = http.StatusConflict
	} else if err.Error() == "invalid status transition" ||
		err.Error() == "device already activated" ||
		err.Error() == "name is required" ||
		err.Error() == "invalid tag key" {
		status = http.StatusBadRequest
	}

	c.JSON(status, gin.H{"code": status, "message": err.Error()})
}

func id(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func page(c *gin.Context, defaultSize int) (int, int) {
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(defaultSize)))
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = defaultSize
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageNumber, pageSize
}

func raw(value string) any {
	if len(value) == 0 {
		return ""
	}

	var legacyString string
	if json.Unmarshal([]byte(value), &legacyString) == nil {
		return legacyString
	}
	return value
}

func deletedAt(value gorm.DeletedAt) *time.Time {
	if value.Valid {
		return &value.Time
	}
	return nil
}
