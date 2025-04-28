package components

import (
	"Dr.uml/backend/Component"
)

type componentsPool struct {
	list []*component.Component
}

func (cp *componentsPool) insert(c *component.Component) {
	cp.list = append(cp.list, c)
}

func (cp *componentsPool) remove(c *component.Component) {
	for i, comp := range cp.list {
		if comp == c {
			cp.list = append(cp.list[:i], cp.list[i+1:]...)
			break
		}
	}
}

func (cp *componentsPool) Len() int {
	return len(cp.list)
}
