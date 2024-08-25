package filters

type FilterBody struct {
	Skip   int         `json:"skip"`
	Take   int         `json:"take"`
	Filter interface{} `json:"filter"`
	Order  []OrderBy   `json:"order"`
}
