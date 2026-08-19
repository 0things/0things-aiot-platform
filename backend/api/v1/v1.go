package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ApiResponse is the unified envelope used across handlers. The payload lives in
// Data, so a single generic type replaces per-endpoint response wrappers.
type ApiResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
} //@name ApiResponse

// Response documents endpoints whose successful result has no payload.
// It preserves the established empty-object response shape.
type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    struct{} `json:"data"`
} //@name ApiSuccessResponse

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	ctx.JSON(http.StatusOK, ApiResponse[any]{
		Code:    errorCodeMap[ErrSuccess],
		Message: ErrSuccess.Error(),
		Data:    data,
	})
}

func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := ApiResponse[any]{Code: errorCodeMap[err], Message: err.Error(), Data: data}
	if _, ok := errorCodeMap[err]; !ok {
		resp = ApiResponse[any]{Code: 500, Message: "unknown error", Data: data}
	}
	ctx.JSON(httpCode, resp)
}

type Error struct {
	Code    int
	Message string
} //@name ApiError

var errorCodeMap = map[error]int{}

func newError(code int, msg string) error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	return err
}
func (e Error) Error() string {
	return e.Message
}
