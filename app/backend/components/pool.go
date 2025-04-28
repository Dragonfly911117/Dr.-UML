package components

import (
	"Dr.uml/backend/component"
)

type Pool struct {
	list []*component.Component
}

func (cp *Pool) insert(c *component.Component) {
	cp.list = append(cp.list, c)
}

func (cp *Pool) remove(c *component.Component) {
	for i, comp := range cp.list {
		if comp == c {
			cp.list = append(cp.list[:i], cp.list[i+1:]...)
			break
		}
	}
}

func (cp *Pool) Len() int {
	return len(cp.list)
}
