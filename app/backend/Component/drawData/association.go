package drawdata

import (
	"Dr.uml/backend/utils"
)

type Association struct {
	AssType    string         `json:"assType"`
	Layer      int            `json:"layer"`
	Start      utils.Point    `json:"start"`
	End        utils.Point    `json:"end"`
	Attributes []AssAttribute `json:"attributes"`
}
