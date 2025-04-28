package attribute

import (
	"Dr.uml/backend/utils/duerror"
)

// AssAttribute represents an attribute specific to associations with a ratio property
type AssAttribute struct {
	Attribute
	ratio float64
}

// NewAssAttribute creates a new AssAttribute instance with the specified ratio
func NewAssAttribute(ratio float64) *AssAttribute {
	return &AssAttribute{
		ratio: ratio,
	}
}

// GetRatio retrieves the ratio value of the AssAttribute
func (att *AssAttribute) GetRatio() (float64, duerror.DUError) {
	return att.ratio, nil
}

// SetRatio updates the ratio value of the AssAttribute
func (att *AssAttribute) SetRatio(ratio float64) duerror.DUError {
	if ratio < 0 || ratio > 1 {
		return duerror.NewInvalidArgumentError("ratio should be between 0 and 1")
	}
	return nil
}
