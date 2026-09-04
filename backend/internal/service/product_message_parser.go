package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	messageParserV1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"github.com/robertkrimen/otto"
)

const defaultProductMessageParserScript = `// 以下为脚本模版，您可以基于以下模版进行脚本编写

/**
 * 将设备自定义topic数据转换为json格式数据, 设备上报数据到物联网平台时调用
 * 入参：topic string 设备上报消息的topic
 * 入参：rawData byte[]数组 不能为空
 * 出参：jsonObj JSON对象 不能为空
 */
function transformPayload(topic, rawData) {
    var jsonObj = {};
    return jsonObj;
}

/**
 * 将设备的自定义格式数据转换为Alink协议的数据，设备上报数据到物联网平台时调用
 * 入参：rawData byte[]数组 不能为空
 * 出参：jsonObj Alink JSON对象 不能为空
 */
function rawDataToProtocol(rawData) {
    var jsonObj = {};
    return jsonObj;
}

/**
 * 将Alink协议的数据转换为设备能识别的格式数据，物联网平台给设备下发数据时调用
 * 入参：jsonObj Alink JSON对象 不能为空
 * 出参：rawData byte[]数组 不能为空
 */
function protocolToRawData(jsonObj) {
    var rawdata = [];
    return rawdata;
}`

const (
	maxMessageParserScriptBytes = 64 * 1024
	maxMessageParserInputBytes  = 4 * 1024
	messageParserTimeout        = 200 * time.Millisecond
)

type ProductMessageParserServiceInterface interface {
	Get(ctx context.Context, productKey string) (*model.ProductMessageParser, bool, error)
	Save(ctx context.Context, productKey, language, script string) (*model.ProductMessageParser, error)
	Execute(ctx context.Context, productKey string, request messageParserV1.ExecuteProductMessageParserRequest) (*messageParserV1.ExecuteProductMessageParserResponse, error)
}

type ProductMessageParserService struct {
	products *repository.ProductRepository
	parsers  *repository.ProductMessageParserRepository
}

func NewProductMessageParserService(products *repository.ProductRepository, parsers *repository.ProductMessageParserRepository) *ProductMessageParserService {
	return &ProductMessageParserService{products: products, parsers: parsers}
}

func (s *ProductMessageParserService) Get(ctx context.Context, productKey string) (*model.ProductMessageParser, bool, error) {
	product, err := s.products.FindByKey(ctx, productKey)
	if err != nil {
		return nil, false, err
	}
	parser, err := s.parsers.FindByProductID(ctx, product.ID)
	if errors.Is(err, repository.ErrNotFound) {
		return &model.ProductMessageParser{ProductID: product.ID, Language: messageParserV1.LanguageJavaScriptES5, Script: defaultProductMessageParserScript}, true, nil
	}
	if err != nil {
		return nil, false, err
	}
	return parser, false, nil
}

func (s *ProductMessageParserService) Save(ctx context.Context, productKey, language, script string) (*model.ProductMessageParser, error) {
	if language != messageParserV1.LanguageJavaScriptES5 {
		return nil, errors.New("unsupported parser language")
	}
	if err := validateMessageParserScript(script); err != nil {
		return nil, err
	}
	product, err := s.products.FindByKey(ctx, productKey)
	if err != nil {
		return nil, err
	}
	parser := &model.ProductMessageParser{ProductID: product.ID, Language: language, Script: script}
	if err := s.parsers.Save(ctx, parser); err != nil {
		return nil, err
	}
	return parser, nil
}

func (s *ProductMessageParserService) Execute(ctx context.Context, productKey string, request messageParserV1.ExecuteProductMessageParserRequest) (*messageParserV1.ExecuteProductMessageParserResponse, error) {
	parser, _, err := s.Get(ctx, productKey)
	if err != nil {
		return nil, err
	}
	if parser.Language != messageParserV1.LanguageJavaScriptES5 {
		return nil, errors.New("unsupported parser language")
	}
	return executeMessageParser(parser.Script, request)
}

func validateMessageParserScript(script string) error {
	if len(script) == 0 || len(script) > maxMessageParserScriptBytes {
		return errors.New("invalid parser script size")
	}
	vm, timedOut, stop := newMessageParserRuntime()
	defer stop()
	if _, err := runMessageParserOperation(timedOut, func() (otto.Value, error) { return vm.Run(script) }); err != nil {
		return fmt.Errorf("invalid parser script: %w", err)
	}
	for _, name := range []string{"transformPayload", "rawDataToProtocol", "protocolToRawData"} {
		value, err := vm.Get(name)
		if err != nil || !value.IsFunction() {
			return fmt.Errorf("parser script must define %s", name)
		}
	}
	return nil
}

func executeMessageParser(script string, request messageParserV1.ExecuteProductMessageParserRequest) (*messageParserV1.ExecuteProductMessageParserResponse, error) {
	if err := validateMessageParserScript(script); err != nil {
		return nil, err
	}
	vm, timedOut, stop := newMessageParserRuntime()
	defer stop()
	if _, err := runMessageParserOperation(timedOut, func() (otto.Value, error) { return vm.Run(script) }); err != nil {
		return nil, fmt.Errorf("invalid parser script: %w", err)
	}

	var functionName string
	var args []interface{}
	switch request.Mode {
	case "device_report":
		raw, err := decodeParserHex(request.RawData)
		if err != nil {
			return nil, err
		}
		functionName = "rawDataToProtocol"
		args = []interface{}{raw}
	case "device_receive":
		raw, err := decodeParserHex(request.RawData)
		if err != nil {
			return nil, err
		}
		functionName = "rawDataToProtocol"
		args = []interface{}{raw}
	case "custom":
		if request.Topic == "" {
			return nil, errors.New("topic is required for custom mode")
		}
		raw, err := decodeParserHex(request.RawData)
		if err != nil {
			return nil, err
		}
		functionName = "transformPayload"
		args = []interface{}{request.Topic, raw}
	default:
		return nil, errors.New("unsupported parser mode")
	}

	value, err := runMessageParserOperation(timedOut, func() (otto.Value, error) {
		return vm.Call(functionName, nil, args...)
	})
	if err != nil {
		return nil, fmt.Errorf("parser execution failed: %w", err)
	}
	if request.Mode == "device_receive" {
		jsonOutput, err := parserJSONObject(value)
		if err != nil {
			return nil, err
		}
		return &messageParserV1.ExecuteProductMessageParserResponse{Mode: request.Mode, JSONOutput: string(jsonOutput)}, nil
	}
	jsonOutput, err := parserJSONObject(value)
	if err != nil {
		return nil, err
	}
	return &messageParserV1.ExecuteProductMessageParserResponse{Mode: request.Mode, JSONOutput: string(jsonOutput)}, nil
}

func newMessageParserRuntime() (*otto.Otto, *atomic.Bool, func()) {
	vm := otto.New()
	vm.Interrupt = make(chan func(), 1)
	timedOut := &atomic.Bool{}
	timer := time.AfterFunc(messageParserTimeout, func() {
		timedOut.Store(true)
		vm.Interrupt <- func() { panic("message parser execution timed out") }
	})
	return vm, timedOut, func() { timer.Stop() }
}

func runMessageParserOperation(timedOut *atomic.Bool, operation func() (otto.Value, error)) (value otto.Value, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if timedOut.Load() {
				err = errors.New("message parser execution timed out")
				return
			}
			panic(recovered)
		}
	}()
	return operation()
}

func decodeParserHex(value string) ([]int, error) {
	clean := strings.TrimPrefix(strings.TrimPrefix(strings.TrimSpace(value), "0x"), "0X")
	clean = strings.ReplaceAll(clean, " ", "")
	if len(clean) == 0 || len(clean)%2 != 0 || len(clean)/2 > maxMessageParserInputBytes {
		return nil, errors.New("invalid hexadecimal input")
	}
	raw, err := hex.DecodeString(clean)
	if err != nil {
		return nil, errors.New("invalid hexadecimal input")
	}
	result := make([]int, len(raw))
	for i, item := range raw {
		result[i] = int(item)
	}
	return result, nil
}

func parserJSONObject(value otto.Value) ([]byte, error) {
	object := value.Object()
	if object == nil || object.Class() != "Object" {
		return nil, errors.New("parser must return a JSON object")
	}
	return json.Marshal(value)
}

func parserByteArray(value otto.Value) ([]byte, error) {
	object := value.Object()
	if object == nil || object.Class() != "Array" {
		return nil, errors.New("parser must return a byte array")
	}
	lengthValue, err := object.Get("length")
	if err != nil {
		return nil, err
	}
	length, err := lengthValue.ToInteger()
	if err != nil || length < 0 || length > maxMessageParserInputBytes {
		return nil, errors.New("parser must return a byte array")
	}
	result := make([]byte, length)
	for i := int64(0); i < length; i++ {
		item, err := object.Get(fmt.Sprint(i))
		if err != nil {
			return nil, err
		}
		number, err := item.ToInteger()
		if err != nil || !item.IsNumber() || number < 0 || number > 255 {
			return nil, errors.New("parser byte array must contain integers from 0 to 255")
		}
		result[i] = byte(number)
	}
	return result, nil
}
