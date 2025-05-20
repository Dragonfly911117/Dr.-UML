package utils

type fileType int

const (
	FiletypeDiagram    fileType = 1 << iota
	FiletypeSubmodule  fileType = 1 << iota
	SupportedFiletypes          = FiletypeDiagram | FiletypeSubmodule
)

type SavedAtt struct {
	Content  string  `json:"content"`
	Size     int     `json:"size"`
	Style    int     `json:"style"`
	FontFile string  `json:"fontFile"`
	Ratio    float64 `json:"ratio,omitempty"`
}

type SavedGad struct {
	GadgetType int        `json:"GadgetType"`
	Point      string     `json:"point"`
	Layer      int        `json:"layer"`
	Attributes []SavedAtt `json:"attributes"`
}

type SavedAss struct {
	AssType    int      `json:"assType"`
	Layer      int      `json:"layer"`
	Parents    []int    `json:"parents"`
	Attributes SavedAtt `json:"attributes"`
}

type SavedFile struct {
	Filetype     fileType   `json:"filetype"`
	LastEdit     string     `json:"lastEdit"`
	Gadgets      []SavedGad `json:"Gadgets"`
	Associations []SavedAss `json:"Associations"`
}
