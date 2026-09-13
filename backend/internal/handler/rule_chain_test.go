package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestBindRuleChainRequestRejectsBlankName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	for _, body := range []string{`{"name":"   ","graph":{"nodes":[{}]}}`, `{"name":"\t\n","graph":{"nodes":[{}]}}`} {
		c.Request = httptest.NewRequest("POST", "/rule-chains", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		_, err := bindRuleChainRequest(c)
		require.EqualError(t, err, "rule chain name is required")
	}
}

func TestBindRuleChainRequestTrimsName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/rule-chains", strings.NewReader(`{"name":"  Temperature alert  ","graph":{"nodes":[{}]}}`))
	c.Request.Header.Set("Content-Type", "application/json")

	req, err := bindRuleChainRequest(c)
	require.NoError(t, err)
	require.Equal(t, "Temperature alert", req.Name)
}
