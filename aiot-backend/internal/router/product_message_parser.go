package router

import "github.com/gin-gonic/gin"

func InitProductMessageParserRouter(deps RouterDeps, r *gin.RouterGroup) {
	parser := deps.ProductMessageParserHandler
	r.GET("/products/key/:productKey/message-parser", parser.Get)
	r.PUT("/products/key/:productKey/message-parser", parser.Put)
	r.POST("/products/key/:productKey/message-parser/execute", parser.Execute)
}
