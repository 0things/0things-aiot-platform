package handler

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func id(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
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
