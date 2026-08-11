package router

import "github.com/gin-gonic/gin"

func InitRuleRouter(deps RouterDeps, r *gin.RouterGroup) {
	rule := deps.RuleHandler
	r.GET("/rules", rule.ListRules)
	r.GET("/rules/:id", rule.GetRule)
	r.POST("/rules", rule.CreateRule)
	r.PUT("/rules/:id", rule.UpdateRule)
	r.DELETE("/rules/:id", rule.DeleteRule)
	r.POST("/rules/:id", rule.RuleAction)
	r.GET("/rules/:id/executions", rule.RuleExecutions)
	r.GET("/rules/available-fields", rule.AvailableFields)
}
