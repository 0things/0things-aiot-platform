package router

import "github.com/gin-gonic/gin"

func InitRuleChainRouter(deps RouterDeps, r *gin.RouterGroup) {
	h := deps.RuleChainHandler
	r.GET("/rule-chains", h.ListRuleChains)
	r.POST("/rule-chains", h.CreateRuleChain)
	r.GET("/rule-chains/:uuid", h.GetRuleChain)
	r.PUT("/rule-chains/:uuid", h.UpdateRuleChain)
	r.DELETE("/rule-chains/:uuid", h.DeleteRuleChain)
}
