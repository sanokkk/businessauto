package filters

type ProductFilter struct {
	CategoryFilter *CategoryRangeFilter `json:"categoryFilter,omitempty"`
	TitleFilter    *TitleFilter         `json:"titleFilter,omitempty"`
	PriceFilter    *PriceFilter         `json:"priceFilter,omitempty"`
	MakerFilter    *MakerFilter         `json:"makerFilter,omitempty"`
}

type CategoryRangeFilter struct {
	Categories []string `json:"categories"`
}

type TitleFilter struct {
	Title string `json:"title"`
}

type PriceFilter struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type MakerFilter struct {
	Makers []string `json:"makers"`
}
