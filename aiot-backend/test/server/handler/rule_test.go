package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"0things-backend/internal/handler"
	"0things-backend/internal/model"
	mock_service "0things-backend/test/mocks/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRuleRouter(mockService *mock_service.MockRuleServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	ruleHandler := handler.NewRuleHandler(h, mockService)

	router.GET("/rules", ruleHandler.ListRules)
	router.GET("/rules/:id", ruleHandler.GetRule)
	router.POST("/rules", ruleHandler.CreateRule)
	router.PUT("/rules/:id", ruleHandler.UpdateRule)
	router.DELETE("/rules/:id", ruleHandler.DeleteRule)
	router.POST("/rules/:id", ruleHandler.RuleAction)
	router.GET("/rules/:id/executions", ruleHandler.RuleExecutions)
	router.POST("/rules/:id/evaluate", ruleHandler.EvaluateRule)

	return router
}

func TestRuleHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockRuleServiceInterface(ctrl)
	router := setupRuleRouter(mockService)

	rules := []model.Rule{{ID: 1, Name: "Temperature Alert"}}
	mockService.EXPECT().List(gomock.Any(), 1, 20, "", "", "").Return(rules, int64(1), nil)

	req, _ := http.NewRequest("GET", "/rules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockRuleServiceInterface(ctrl)
	router := setupRuleRouter(mockService)

	rule := &model.Rule{ID: 1, Name: "Temperature Alert"}
	mockService.EXPECT().Get(gomock.Any(), int64(1)).Return(rule, nil)

	req, _ := http.NewRequest("GET", "/rules/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockRuleServiceInterface(ctrl)
	router := setupRuleRouter(mockService)

	mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]string{"name": "Temperature Alert"})
	req, _ := http.NewRequest("POST", "/rules", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockRuleServiceInterface(ctrl)
	router := setupRuleRouter(mockService)

	mockService.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/rules/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_SetStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockRuleServiceInterface(ctrl)
	router := setupRuleRouter(mockService)

	rule := &model.Rule{ID: 1, Name: "Temperature Alert", Status: "enabled"}
	mockService.EXPECT().SetStatus(gomock.Any(), int64(1), "enabled").Return(rule, nil)

	req, _ := http.NewRequest("POST", "/rules/1:enable", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
