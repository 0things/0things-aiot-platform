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

func setupRuleRouterFull(mockService *mock_service.MockRuleServiceInterface) *gin.Engine {
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
	router.GET("/rules/available-fields", ruleHandler.AvailableFields)

	return router
}

func TestRuleHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockRuleServiceInterface(ctrl)
	router := setupRuleRouterFull(mockService)

	existingRule := &model.Rule{ID: 1, Name: "Temperature Alert"}
	mockService.EXPECT().Get(gomock.Any(), int64(1)).Return(existingRule, nil)
	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "Temperature Alert",
		"description": "Alert when temperature exceeds threshold",
	})
	req, _ := http.NewRequest("PUT", "/rules/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_RuleExecutions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockRuleServiceInterface(ctrl)
	router := setupRuleRouterFull(mockService)

	executions := []model.RuleExecution{{ID: 1, RuleID: 1, RuleName: "Temperature Alert"}}
	mockService.EXPECT().ListExecutions(gomock.Any(), int64(1), 1, 20).Return(executions, int64(1), nil)

	req, _ := http.NewRequest("GET", "/rules/1/executions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_EvaluateRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockRuleServiceInterface(ctrl)
	router := setupRuleRouterFull(mockService)

	execution := &model.RuleExecution{ID: 1, RuleID: 1, Status: "success"}
	mockService.EXPECT().Evaluate(gomock.Any(), int64(1)).Return(execution, nil)

	req, _ := http.NewRequest("POST", "/rules/1/evaluate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_AvailableFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockRuleServiceInterface(ctrl)
	router := setupRuleRouterFull(mockService)

	// AvailableFields doesn't use the service, it returns static data
	req, _ := http.NewRequest("GET", "/rules/available-fields", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
