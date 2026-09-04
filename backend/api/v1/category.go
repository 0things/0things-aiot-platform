package v1

type Category struct {
	ID       int64      `json:"id"`
	ParentID *int64     `json:"parentId,omitempty"`
	Name     string     `json:"name"`
	Sort     int        `json:"sort"`
	Enabled  bool       `json:"enabled"`
	Children []Category `json:"children,omitempty"`
}
