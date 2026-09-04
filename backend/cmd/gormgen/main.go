//go:generate go run .

package main

import (
	"aiot-backend/internal/model"
	"gorm.io/gen"
	"path/filepath"
	"runtime"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	g := gen.NewGenerator(gen.Config{
		OutPath: filepath.Join(filepath.Dir(file), "../../internal/dal/query"),
	})

	g.ApplyBasic(
		model.Product{},
		model.Device{},
		model.DeviceState{},
		model.DeviceTag{},
		model.DeviceShadow{},
		model.DeviceShadowHistory{},
		model.DeviceEvent{},
		model.DeviceServiceInvocation{},
		model.OTAPackage{},
		model.UpgradeBatch{},
		model.DeviceUpgradeStatus{},
		model.ProductTSL{},
		model.ProductMessageParser{},
		model.DevicePushRecord{},
		model.DeviceGroup{},
		model.DeviceGroupMember{},
		model.SceneLinkage{},
		model.SceneLinkageDetail{},
	)

	g.Execute()
}
