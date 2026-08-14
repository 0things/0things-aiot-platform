//go:generate go run .

package main

import (
	"0things-backend/internal/model"
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
		model.OTAPackage{},
		model.UpgradeBatch{},
		model.DeviceUpgradeStatus{},
		model.ProductTSL{},
		model.DevicePushRecord{},
		model.Rule{},
		model.RuleExecution{},
		model.Alert{},
	)

	g.Execute()
}
