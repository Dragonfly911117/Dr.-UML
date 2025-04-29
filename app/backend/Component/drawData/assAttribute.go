package drawdata

type AssAttribute struct {
	Attribute Attribute `json:"attribute"`
	Content   string    `json:"content"`
	Size      int       `json:"size"`
	Style     int       `json:"style"`
	Ratio     float64   `json:"ratio"`
}
