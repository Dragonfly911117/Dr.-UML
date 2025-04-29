package component

import (
	"Dr.uml/backend/component/attribute"
	"Dr.uml/backend/utils/duerror"
)

type AssociationType int

const (
	Extension      = 1 << iota // 0x01
	Implementation = 1 << iota // 0x02
	Composition    = 1 << iota // 0x04
	Dependency     = 1 << iota // 0x08
	supportedType  = Extension | Implementation | Composition | Dependency
)

type Association struct {
	assType    AssociationType
	layer      int
	attributes []*attribute.AssAttribute
	parents    [2]*Gadget
}

func NewAssociation(parents [2]*Gadget, assType AssociationType) (*Association, duerror.DUError) {
	if assType&supportedType == 0 {
		return nil, duerror.NewInvalidArgumentError("unsupported association type")
	}
	if parents[0] == nil || parents[1] == nil {
		return nil, duerror.NewInvalidArgumentError("parents are nil")
	}
	return &Association{
		parents: [2]*Gadget{parents[0], parents[1]},
	}, nil
}

// MoveAttribute returns an Invalid argument error if the index is out of range or if the ratio is not between 0 and 1
func (a *Association) MoveAttribute(index int, ratio float64) duerror.DUError {
	if index < 0 || index >= len(a.attributes) {
		return duerror.NewInvalidArgumentError("index out of range")
	}
	return a.attributes[index].SetRatio(ratio)
}

func (a *Association) GetAssType() AssociationType {
	return a.assType
}

func (a *Association) SetAssType(assType AssociationType) {
	a.assType = assType
}

func (a *Association) GetLayer() int {
	return a.layer
}

func (a *Association) SetLayer(layer int) {
	a.layer = layer
}

// AddAttribute adds an attribute to the association. Returns an error if the attribute is nil.
func (a *Association) AddAttribute(attribute *attribute.AssAttribute) duerror.DUError {
	if attribute == nil {
		return duerror.NewInvalidArgumentError("attribute is nil")
	}
	a.attributes = append(a.attributes, attribute)
	return nil
}

// GetAttributes returns the attributes of the association. Returns an error if no attributes are found.
func (a *Association) GetAttributes() ([]*attribute.AssAttribute, duerror.DUError) {
	if len(a.attributes) == 0 {
		return nil, duerror.NewInvalidArgumentError("no attributes found")
	}
	return a.attributes, nil
}

func (a *Association) RemoveAttribute(index int) duerror.DUError {
	if index < 0 || index >= len(a.attributes) {
		return duerror.NewInvalidArgumentError("index out of range")
	}
	a.attributes = append(a.attributes[:index], a.attributes[index+1:]...)
	return nil
}

func (a *Association) GetParentStart() *Gadget {
	return a.parents[0]
}

func (a *Association) GetParentEnd() *Gadget {
	return a.parents[1]
}

func (a *Association) SetParentStart(gadget *Gadget) duerror.DUError {
	if gadget == nil {
		return duerror.NewInvalidArgumentError("gadget is nil")
	}
	a.parents[0] = gadget
	return nil
}

func (a *Association) SetParentEnd(gadget *Gadget) duerror.DUError {
	if gadget == nil {
		return duerror.NewInvalidArgumentError("gadget is nil")
	}
	a.parents[1] = gadget
	return nil
}
