package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	eventV1 "aiot-backend/api/event/v1"
	messageParserV1 "aiot-backend/api/message_parser/v1"
	productTSLV1 "aiot-backend/api/product_tsl/v1"
	sceneLinkageV1 "aiot-backend/api/scene_linkage/v1"
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/jwt"
	"aiot-backend/pkg/log"
	"aiot-backend/test/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
)

func baseHandler(t *testing.T) *Handler {
	t.Helper()
	conf := viper.New()
	conf.Set("log.log_file_name", filepath.Join(t.TempDir(), "out.log"))
	return NewHandler(log.NewLog(conf))
}

func hctx(method, target string, body []byte, claims *jwt.MyCustomClaims, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
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
	if claims != nil {
		c.Set("claims", claims)
	}
	if params != nil {
		c.Params = params
	}
	return c, w
}

func idParam(v string) gin.Params { return gin.Params{{Key: "id", Value: v}} }

// ---------------- UserHandler ----------------
func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	h := NewUserHandler(baseHandler(t), service.UserService(mock_service.NewMockUserService(ctrl)))
	m := mock_service.NewMockUserService(ctrl)
	h = NewUserHandler(baseHandler(t), m)
	m.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil)
	c, w := hctx(http.MethodPost, "/register", []byte(`{"email":"a@b.com","password":"x"}`), nil, nil)
	h.Register(c)
	if w.Code != http.StatusOK {
		t.Fatalf("got %d", w.Code)
	}
}

func TestUserHandler_Register_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockUserService(ctrl)
	h := NewUserHandler(baseHandler(t), m)
	c, w := hctx(http.MethodPost, "/register", []byte(`{bad`), nil, nil)
	h.Register(c)
	if w.Code == http.StatusOK {
		t.Fatalf("expected error")
	}
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockUserService(ctrl)
	h := NewUserHandler(baseHandler(t), m)
	m.EXPECT().Login(gomock.Any(), gomock.Any()).Return("tok", nil)
	c, w := hctx(http.MethodPost, "/login", []byte(`{"email":"a@b.com","password":"x"}`), nil, nil)
	h.Login(c)
	if w.Code != http.StatusOK {
		t.Fatalf("got %d", w.Code)
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockUserService(ctrl)
	h := NewUserHandler(baseHandler(t), m)
	m.EXPECT().GetProfile(gomock.Any(), "u1").Return(&v1.GetProfileResponseData{UserId: "u1"}, nil)
	c, w := hctx(http.MethodGet, "/user", nil, &jwt.MyCustomClaims{UserId: "u1"}, nil)
	h.GetProfile(c)
	if w.Code != http.StatusOK {
		t.Fatalf("got %d", w.Code)
	}
}

func TestUserHandler_GetProfile_NoClaims(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockUserService(ctrl)
	h := NewUserHandler(baseHandler(t), m)
	c, w := hctx(http.MethodGet, "/user", nil, nil, nil)
	h.GetProfile(c)
	if w.Code == http.StatusOK {
		t.Fatalf("expected unauthorized")
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockUserService(ctrl)
	h := NewUserHandler(baseHandler(t), m)
	m.EXPECT().UpdateProfile(gomock.Any(), "u1", gomock.Any()).Return(nil)
	c, w := hctx(http.MethodPut, "/user", []byte(`{"email":"a@b.com","nickname":"n"}`), &jwt.MyCustomClaims{UserId: "u1"}, nil)
	h.UpdateProfile(c)
	if w.Code != http.StatusOK {
		t.Fatalf("got %d", w.Code)
	}
}

// ---------------- ProductHandler ----------------
func TestProductHandler_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockProductServiceInterface(ctrl)
	h := NewProductHandler(baseHandler(t), m)

	m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Product{ID: 1}, nil)
	c, w := hctx(http.MethodPost, "/products", []byte(`{"name":"p"}`), nil, nil)
	h.Create(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Create got %d", w.Code)
	}

	m.EXPECT().Get(gomock.Any(), int64(1)).Return(&model.Product{ID: 1}, nil)
	c, w = hctx(http.MethodGet, "/products/1", nil, nil, idParam("1"))
	h.Get(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Get got %d", w.Code)
	}

	m.EXPECT().GetByKey(gomock.Any(), "P001").Return(&model.Product{ID: 1}, nil)
	c, w = hctx(http.MethodGet, "/products/key/P001", nil, nil, gin.Params{{Key: "productKey", Value: "P001"}})
	h.GetByKey(c)
	if w.Code != http.StatusOK {
		t.Fatalf("GetByKey got %d", w.Code)
	}

	m.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.Product{}, int64(0), nil)
	c, w = hctx(http.MethodGet, "/products", nil, nil, nil)
	h.List(c)
	if w.Code != http.StatusOK {
		t.Fatalf("List got %d", w.Code)
	}

	m.EXPECT().GetByKey(gomock.Any(), "P001").Return(&model.Product{ID: 1}, nil)
	m.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
	c, w = hctx(http.MethodPut, "/products/key/P001", []byte(`{"name":"x"}`), nil, gin.Params{{Key: "productKey", Value: "P001"}})
	h.Update(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Update got %d", w.Code)
	}

	m.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)
	c, w = hctx(http.MethodDelete, "/products/1", nil, nil, idParam("1"))
	h.Delete(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Delete got %d", w.Code)
	}

	m.EXPECT().Restore(gomock.Any(), int64(1)).Return(&model.Product{ID: 1}, nil)
	c, w = hctx(http.MethodPost, "/products/1/restore", nil, nil, idParam("1"))
	h.Restore(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Restore got %d", w.Code)
	}
}

func TestProductHandler_Get_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockProductServiceInterface(ctrl)
	h := NewProductHandler(baseHandler(t), m)
	c, w := hctx(http.MethodGet, "/products/abc", nil, nil, idParam("abc"))
	h.Get(c)
	if w.Code == http.StatusOK {
		t.Fatalf("expected error")
	}
}

// ---------------- DeviceEventHandler ----------------
func TestDeviceEventHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	h := NewDeviceEventHandler(baseHandler(t), m)
	m.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]model.DeviceEvent{}, int64(0), nil)
	c, w := hctx(http.MethodGet, "/device-events", nil, nil, nil)
	h.ListDeviceEvents(c)
	if w.Code != http.StatusOK {
		t.Fatalf("got %d", w.Code)
	}
}

func TestDeviceEventHandler_List_BadTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	h := NewDeviceEventHandler(baseHandler(t), m)
	c, w := hctx(http.MethodGet, "/device-events?start_at=not-a-time", nil, nil, nil)
	h.ListDeviceEvents(c)
	if w.Code == http.StatusOK {
		t.Fatalf("expected bad request")
	}
	_ = eventV1.DeviceEvent{}
}

// ---------------- ProductTSLHandler ----------------
func TestProductTSLHandler_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockProductTSLServiceInterface(ctrl)
	h := NewProductTSLHandler(baseHandler(t), m)

	m.EXPECT().Get(gomock.Any(), "1").Return(&model.ProductTSL{}, nil)
	c, w := hctx(http.MethodGet, "/products/1/tsl", nil, nil, idParam("1"))
	h.Get(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Get got %d", w.Code)
	}

	m.EXPECT().Upsert(gomock.Any(), "1", gomock.Any()).Return(nil)
	c, w = hctx(http.MethodPut, "/products/1/tsl", []byte(`{"tsl":"x"}`), nil, idParam("1"))
	h.Put(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Put got %d", w.Code)
	}

	m.EXPECT().Delete(gomock.Any(), "1").Return(nil)
	c, w = hctx(http.MethodDelete, "/products/1/tsl", nil, nil, idParam("1"))
	h.Delete(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Delete got %d", w.Code)
	}
	_ = productTSLV1.ProductTSL{}
}

// ---------------- ProductMessageParserHandler ----------------
func TestProductMessageParserHandler_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockProductMessageParserServiceInterface(ctrl)
	h := NewProductMessageParserHandler(baseHandler(t), m)

	m.EXPECT().Get(gomock.Any(), "P001").Return(&model.ProductMessageParser{}, false, nil)
	c, w := hctx(http.MethodGet, "/products/key/P001/message-parser", nil, nil, gin.Params{{Key: "productKey", Value: "P001"}})
	h.Get(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Get got %d", w.Code)
	}

	m.EXPECT().Save(gomock.Any(), "P001", gomock.Any(), gomock.Any()).Return(&model.ProductMessageParser{}, nil)
	c, w = hctx(http.MethodPut, "/products/key/P001/message-parser", []byte(`{"language":"javascript-es5","script":"x"}`), nil, gin.Params{{Key: "productKey", Value: "P001"}})
	h.Put(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Put got %d", w.Code)
	}

	m.EXPECT().Execute(gomock.Any(), "P001", gomock.Any()).Return(&messageParserV1.ExecuteProductMessageParserResponse{}, nil)
	c, w = hctx(http.MethodPost, "/products/key/P001/message-parser/execute", []byte(`{"mode":"custom","topic":"t","rawData":"00"}`), nil, gin.Params{{Key: "productKey", Value: "P001"}})
	h.Execute(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Execute got %d", w.Code)
	}
}

// ---------------- SceneLinkageHandler ----------------
func TestSceneLinkageHandler_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockSceneLinkageServiceInterface(ctrl)
	h := NewSceneLinkageHandler(baseHandler(t), m)

	m.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.SceneLinkage{}, int64(0), nil)
	c, w := hctx(http.MethodGet, "/scene-linkages", nil, nil, nil)
	h.ListSceneLinkages(c)
	if w.Code != http.StatusOK {
		t.Fatalf("List got %d", w.Code)
	}

	m.EXPECT().Get(gomock.Any(), int64(1)).Return(&model.SceneLinkage{ID: 1}, nil)
	c, w = hctx(http.MethodGet, "/scene-linkages/1", nil, nil, idParam("1"))
	h.GetSceneLinkage(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Get got %d", w.Code)
	}

	m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	c, w = hctx(http.MethodPost, "/scene-linkages", []byte(`{"name":"s","enable":1}`), nil, nil)
	h.CreateSceneLinkage(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Create got %d", w.Code)
	}

	m.EXPECT().Get(gomock.Any(), int64(1)).Return(&model.SceneLinkage{ID: 1}, nil)
	m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	c, w = hctx(http.MethodPut, "/scene-linkages/1", []byte(`{"name":"s"}`), nil, idParam("1"))
	h.UpdateSceneLinkage(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Update got %d", w.Code)
	}

	m.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)
	c, w = hctx(http.MethodDelete, "/scene-linkages/1", nil, nil, idParam("1"))
	h.DeleteSceneLinkage(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Delete got %d", w.Code)
	}
	_ = sceneLinkageV1.SceneLinkage{}
}

// ---------------- SceneLinkageDetailHandler ----------------
func TestSceneLinkageDetailHandler_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_service.NewMockSceneLinkageDetailServiceInterface(ctrl)
	h := NewSceneLinkageDetailHandler(baseHandler(t), m)

	m.EXPECT().GetBySceneID(gomock.Any(), int64(1)).Return(&model.SceneLinkageDetail{SceneID: 1}, nil)
	c, w := hctx(http.MethodGet, "/scene-linkages/1/detail", nil, nil, idParam("1"))
	h.GetSceneLinkageDetail(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Get got %d", w.Code)
	}

	m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	c, w = hctx(http.MethodPost, "/scene-linkages/1/detail", []byte(`{"triggerConfig":{},"actionConfig":{}}`), nil, idParam("1"))
	h.CreateSceneLinkageDetail(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Create got %d", w.Code)
	}

	m.EXPECT().GetBySceneID(gomock.Any(), int64(1)).Return(&model.SceneLinkageDetail{SceneID: 1}, nil)
	m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	c, w = hctx(http.MethodPut, "/scene-linkages/1/detail", []byte(`{"triggerConfig":{},"actionConfig":{}}`), nil, idParam("1"))
	h.UpdateSceneLinkageDetail(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Update got %d", w.Code)
	}
}
