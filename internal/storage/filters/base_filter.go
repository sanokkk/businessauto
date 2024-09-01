package filters

type FilterBody struct {
	Skip   *int        `json:"skip" binding:"required"`
	Take   *int        `json:"take" binding:"required"`
	Filter interface{} `json:"filter"`
	Order  []OrderBy   `json:"order" binding:"required"`
}
