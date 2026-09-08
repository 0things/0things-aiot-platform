package server

import (
	"transport-http/internal/router"
	"transport-http/pkg/log"
	pkgHTTP "transport-http/pkg/server/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewHTTPServer(
	conf *viper.Viper,
	logger *log.Logger,
	r *router.Router,
) *pkgHTTP.Server {
	if conf.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	r.Register(engine)

	s := pkgHTTP.NewServer(
		engine,
		logger,
		pkgHTTP.WithServerHost(conf.GetString("http.host")),
		pkgHTTP.WithServerPort(conf.GetInt("http.port")),
	)
	return s
}
