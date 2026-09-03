package v1

// PageRequest contains the common pagination query parameters.
type PageRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`             // Page number (1-based, default 1)
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100"` // Page size (1-100, default 10)
} //@name PageRequest
