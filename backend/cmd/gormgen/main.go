//go:generate go run .

package main

import (
	"aiot-backend/internal/model"
	"path/filepath"
	"runtime"

	"gorm.io/gen"
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
		model.OTAUpgradeBatch{},
		model.OTADeviceUpgradeStatus{},
		model.ProductTSL{},
		model.ProductMessageParser{},
		model.DevicePushRecord{},
		model.DeviceGroup{},
		model.DeviceGroupMember{},
	)

	g.Execute()
}
