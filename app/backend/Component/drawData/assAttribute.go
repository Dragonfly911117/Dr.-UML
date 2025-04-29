package drawdata

type AssAttribute struct {
	Content   string  `json:"content"`
	FontSize  int     `json:"fontSize"`
	FontStyle int     `json:"fontStyle"`
	FontFile  string  `json:"fontFile"`
	Ratio     float64 `json:"ratio"`
}
