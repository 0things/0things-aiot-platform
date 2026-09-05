package dto

// ListDeviceGroupDevicesQuery defines query parameters for filtering and paginating device group devices.
type ListDeviceGroupDevicesQuery struct {
	GroupUUID   string
	GroupID     int64
	GroupType   string
	Rule        string
	Page        int
	PageSize    int
	ProductKeys []string
	Search      string
}
