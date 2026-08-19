package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"aiot-backend/internal/handler"
	"aiot-backend/pkg/config"
	"aiot-backend/pkg/log"
	v1 "aiot-backend/api/v1"
	"aiot-backend/pkg/jwt"
	mock_service "aiot-backend/test/mocks/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var testLogger *log.Logger

func init() {
	os.Setenv("APP_CONF", "../../../config/local.yml")
	conf := config.NewConfig("config/local.yml")
	testLogger = log.NewLog(conf)
}

func setupUserRouter(t *testing.T) (*mock_service.MockUserService, *gin.Engine) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockService := mock_service.NewMockUserService(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHandler(testLogger)
	uh := handler.NewUserHandler(h, mockService)

	router.POST("/register", uh.Register)
	router.POST("/login", uh.Login)

	return mockService, router
}

func TestUserHandler_Register_BadJSON(t *testing.T) {
	_, router := setupUserRouter(t)

	req, _ := http.NewRequest("POST", "/register", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Register_ServiceError(t *testing.T) {
	mockService, router := setupUserRouter(t)
	params := v1.RegisterRequest{Email: "test@test.com", Password: "pass"}
	mockService.EXPECT().Register(gomock.Any(), &params).Return(errors.New("db error"))

	body, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserHandler_Register_Success(t *testing.T) {
	mockService, router := setupUserRouter(t)
	params := v1.RegisterRequest{Email: "test@test.com", Password: "pass"}
	mockService.EXPECT().Register(gomock.Any(), &params).Return(nil)

	body, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_Login_BadJSON(t *testing.T) {
	_, router := setupUserRouter(t)

	req, _ := http.NewRequest("POST", "/login", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Login_ServiceError(t *testing.T) {
	mockService, router := setupUserRouter(t)
	params := v1.LoginRequest{Email: "test@test.com", Password: "wrong"}
	mockService.EXPECT().Login(gomock.Any(), &params).Return("", errors.New("unauthorized"))

	body, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_Login_Success(t *testing.T) {
	mockService, router := setupUserRouter(t)
	params := v1.LoginRequest{Email: "test@test.com", Password: "pass"}
	mockService.EXPECT().Login(gomock.Any(), &params).Return("token123", nil)

	body, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "token123", data["accessToken"])
}

func TestUserHandler_GetProfile_NoUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	_ = mock_service.NewMockUserService(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHandler(testLogger)
	uh := handler.NewUserHandler(h, nil)
	router.GET("/profile", uh.GetProfile)

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_GetProfile_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockService := mock_service.NewMockUserService(ctrl)
	mockService.EXPECT().GetProfile(gomock.Any(), "user123").Return(nil, errors.New("not found"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHandler(testLogger)
	uh := handler.NewUserHandler(h, mockService)

	router.Use(func(c *gin.Context) {
		c.Set("claims", &jwt.MyCustomClaims{UserId: "user123"})
		c.Next()
	})
	router.GET("/profile", uh.GetProfile)

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockService := mock_service.NewMockUserService(ctrl)
	mockService.EXPECT().GetProfile(gomock.Any(), "user123").Return(&v1.GetProfileResponseData{UserId: "user123", Nickname: "test"}, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHandler(testLogger)
	uh := handler.NewUserHandler(h, mockService)

	router.Use(func(c *gin.Context) {
		c.Set("claims", &jwt.MyCustomClaims{UserId: "user123"})
		c.Next()
	})
	router.GET("/profile", uh.GetProfile)

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_UpdateProfile_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockService := mock_service.NewMockUserService(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHandler(testLogger)
	uh := handler.NewUserHandler(h, mockService)

	router.Use(func(c *gin.Context) {
		c.Set("claims", &jwt.MyCustomClaims{UserId: "user123"})
		c.Next()
	})
	router.PUT("/profile", uh.UpdateProfile)

	req, _ := http.NewRequest("PUT", "/profile", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UpdateProfile_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockService := mock_service.NewMockUserService(ctrl)
	params := v1.UpdateProfileRequest{Email: "new@test.com", Nickname: "new"}
	mockService.EXPECT().UpdateProfile(gomock.Any(), "user123", &params).Return(errors.New("db error"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHandler(testLogger)
	uh := handler.NewUserHandler(h, mockService)

	router.Use(func(c *gin.Context) {
		c.Set("claims", &jwt.MyCustomClaims{UserId: "user123"})
		c.Next()
	})
	router.PUT("/profile", uh.UpdateProfile)

	body, _ := json.Marshal(params)
	req, _ := http.NewRequest("PUT", "/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserHandler_UpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockService := mock_service.NewMockUserService(ctrl)
	params := v1.UpdateProfileRequest{Email: "new@test.com", Nickname: "new"}
	mockService.EXPECT().UpdateProfile(gomock.Any(), "user123", &params).Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := handler.NewHandler(testLogger)
	uh := handler.NewUserHandler(h, mockService)

	router.Use(func(c *gin.Context) {
		c.Set("claims", &jwt.MyCustomClaims{UserId: "user123"})
		c.Next()
	})
	router.PUT("/profile", uh.UpdateProfile)

	body, _ := json.Marshal(params)
	req, _ := http.NewRequest("PUT", "/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserIdFromCtx_NoClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	result := handler.GetUserIdFromCtx(c)
	assert.Equal(t, "", result)
}

func TestGetUserIdFromCtx_WithClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	c.Set("claims", &jwt.MyCustomClaims{UserId: "user123"})
	result := handler.GetUserIdFromCtx(c)
	assert.Equal(t, "user123", result)
}
