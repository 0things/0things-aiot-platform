package router

import "github.com/gin-gonic/gin"

func InitRuleNodeDefinitionRouter(deps RouterDeps, r *gin.RouterGroup) {
	r.GET("/rule-node-definitions", deps.RuleNodeDefinitionHandler.ListRuleNodeDefinitions)
}
