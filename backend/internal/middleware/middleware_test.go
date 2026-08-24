package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"aiot-backend/pkg/jwt"
	"aiot-backend/pkg/log"

	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func mctx(method, target string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var r *http.Request
	if body != nil {
		r = httptest.NewRequest(method, target, bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
	} else {
		r = httptest.NewRequest(method, target, nil)
	}
	c.Request = r
	return c, w
}

func mlogger() *log.Logger {
	conf := viper.New()
	conf.Set("log.log_file_name", "/dev/null")
	return log.NewLog(conf)
}

func TestCORSMiddleware_Options(t *testing.T) {
	h := CORSMiddleware()
	c, w := mctx(http.MethodOptions, "/", nil)
	h(c)
	if !c.IsAborted() || w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 abort, got %d aborted=%v", w.Code, c.IsAborted())
	}
}

func TestCORSMiddleware_Get(t *testing.T) {
	h := CORSMiddleware()
	c, _ := mctx(http.MethodGet, "/", nil)
	h(c)
	if c.IsAborted() {
		t.Fatalf("GET should not abort")
	}
}

func TestSignMiddleware_MissingHeader(t *testing.T) {
	conf := viper.New()
	conf.Set("security.api_sign.app_key", "k")
	conf.Set("security.api_sign.app_security", "s")
	h := SignMiddleware(mlogger(), conf)
	c, _ := mctx(http.MethodGet, "/", nil)
	h(c)
	if !c.IsAborted() {
		t.Fatalf("expected abort on missing header")
	}
}

func TestSignMiddleware_Valid(t *testing.T) {
	conf := viper.New()
	conf.Set("security.api_sign.app_key", "k")
	conf.Set("security.api_sign.app_security", "s")
	h := SignMiddleware(mlogger(), conf)
	c, _ := mctx(http.MethodGet, "/", nil)
	ts, nonce, ver := "123", "abc", "v1"
	data := map[string]string{"AppKey": "k", "Timestamp": ts, "Nonce": nonce, "AppVersion": ver}
	keys := []string{"AppKey", "Timestamp", "Nonce", "AppVersion"}
	sort.Slice(keys, func(i, j int) bool { return strings.ToLower(keys[i]) < strings.ToLower(keys[j]) })
	var str string
	for _, k := range keys {
		str += k + data[k]
	}
	str += "s"
	c.Request.Header.Set("Timestamp", ts)
	c.Request.Header.Set("Nonce", nonce)
	c.Request.Header.Set("App-Version", ver)
	c.Request.Header.Set("Sign", strings.ToUpper(cryptor.Md5String(str)))
	h(c)
	if c.IsAborted() {
		t.Fatalf("valid sign should not abort")
	}
}

func TestStrictAuth_NoToken(t *testing.T) {
	conf := viper.New()
	conf.Set("security.jwt.key", "secret")
	j := jwt.NewJwt(conf)
	h := StrictAuth(j, mlogger())
	c, w := mctx(http.MethodGet, "/", nil)
	h(c)
	if !c.IsAborted() || w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 abort, got %d", w.Code)
	}
}

func TestStrictAuth_InvalidToken(t *testing.T) {
	conf := viper.New()
	conf.Set("security.jwt.key", "secret")
	j := jwt.NewJwt(conf)
	h := StrictAuth(j, mlogger())
	c, w := mctx(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "bad.token.value")
	h(c)
	if !c.IsAborted() || w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 abort, got %d", w.Code)
	}
}

func TestStrictAuth_ValidToken(t *testing.T) {
	conf := viper.New()
	conf.Set("security.jwt.key", "secret")
	j := jwt.NewJwt(conf)
	token, err := j.GenToken("u1", 1, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	h := StrictAuth(j, mlogger())
	c, _ := mctx(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", token)
	h(c)
	if c.IsAborted() {
		t.Fatalf("valid token should not abort")
	}
	if _, ok := c.Get("claims"); !ok {
		t.Fatalf("claims should be set")
	}
}

func TestNoStrictAuth_NoToken(t *testing.T) {
	conf := viper.New()
	conf.Set("security.jwt.key", "secret")
	j := jwt.NewJwt(conf)
	h := NoStrictAuth(j, mlogger())
	c, _ := mctx(http.MethodGet, "/", nil)
	h(c)
	if c.IsAborted() {
		t.Fatalf("no token should pass through")
	}
}

func TestNoStrictAuth_InvalidToken(t *testing.T) {
	conf := viper.New()
	conf.Set("security.jwt.key", "secret")
	j := jwt.NewJwt(conf)
	h := NoStrictAuth(j, mlogger())
	c, _ := mctx(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "bad")
	h(c)
	if c.IsAborted() {
		t.Fatalf("invalid token should pass through")
	}
}

func TestNoStrictAuth_ValidToken(t *testing.T) {
	conf := viper.New()
	conf.Set("security.jwt.key", "secret")
	j := jwt.NewJwt(conf)
	token, _ := j.GenToken("u1", 1, time.Now().Add(time.Hour))
	h := NoStrictAuth(j, mlogger())
	c, _ := mctx(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", token)
	h(c)
	if c.IsAborted() {
		t.Fatalf("valid token should pass through")
	}
	if _, ok := c.Get("claims"); !ok {
		t.Fatalf("claims should be set")
	}
}

func TestRequestLogMiddleware(t *testing.T) {
	h := RequestLogMiddleware(mlogger())
	c, _ := mctx(http.MethodPost, "/x", []byte(`{"a":1}`))
	h(c)
	if c.IsAborted() {
		t.Fatalf("should not abort")
	}
}

func TestResponseLogMiddleware(t *testing.T) {
	h := ResponseLogMiddleware(mlogger())
	c, w := mctx(http.MethodGet, "/x", nil)
	c.Status(http.StatusOK)
	h(c)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected code %d", w.Code)
	}
}
