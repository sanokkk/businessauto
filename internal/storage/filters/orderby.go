package filters

type OrderBy struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc"`
}
