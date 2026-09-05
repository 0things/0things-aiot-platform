package dto

// ListDevicesQuery defines query parameters for filtering and paginating devices.
type ListDevicesQuery struct {
	Page      int
	PageSize  int
	ProductID int64
	States    []string
	Enabled   *bool
	Search    string
}

