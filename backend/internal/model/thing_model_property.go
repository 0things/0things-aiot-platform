package model

import "time"

// ThingModelProperty is a TSL-defined property with its most recent telemetry point.
type ThingModelProperty struct {
	Identifier string
	Name       string
	DataType   string
	Unit       string
	AccessMode string
	Value      any
	ReportedAt *time.Time
}
