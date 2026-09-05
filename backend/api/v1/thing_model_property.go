package v1

import "time"

type ThingModelProperty struct {
	Identifier string     `json:"identifier"`
	Name       string     `json:"name"`
	DataType   string     `json:"dataType"`
	Unit       string     `json:"unit"`
	AccessMode string     `json:"accessMode"`
	Value      any        `json:"value" swaggertype:"object"`
	ReportedAt *time.Time `json:"reportedAt"`
} //@name ThingModelProperty

type GetDeviceThingModelPropertiesResponse struct {
	Items []ThingModelProperty `json:"items"`
} //@name GetDeviceThingModelPropertiesResponse
