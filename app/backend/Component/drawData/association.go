package drawdata

import (
	"Dr.uml/backend/component"
	"Dr.uml/backend/component/attribute"
	"Dr.uml/backend/utils"
)

type association struct {
	AssType    component.AssociationType `json:"assType"`
	Layer      int                       `json:"layer"`
	Start      utils.Point               `json:"start"`
	End        utils.Point               `json:"end"`
	Attributes []attribute.AssAttribute  `json:"attributes"`
}
