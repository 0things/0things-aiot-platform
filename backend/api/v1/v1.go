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
	code := codeForError(err)
	ctx.JSON(httpCode, ApiResponse[any]{Code: code, Message: err.Error(), Data: data})
}

// codeForError resolves the registered HTTP code for a sentinel error. Some
// errors (e.g. validator.ValidationErrors) are not hashable, so the map lookup
// is guarded to avoid a panic on invalid request payloads.
func codeForError(err error) int {
	defer func() { _ = recover() }()
	if code, ok := errorCodeMap[err]; ok {
		return code
	}
	return 500
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
