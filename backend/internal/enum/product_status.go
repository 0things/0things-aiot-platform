package enum

// ProductStatus describes the lifecycle state of a product.
type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"   // 产品可正常使用。
	ProductStatusInactive ProductStatus = "inactive" // 产品已暂停使用。
	ProductStatusArchived ProductStatus = "archived" // 产品已归档并停止使用。
)
