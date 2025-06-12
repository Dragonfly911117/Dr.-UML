// VIBE CODING

package umldiagram

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"Dr.uml/backend/component/attribute"

	"Dr.uml/backend/component"
	"Dr.uml/backend/drawdata"
	"Dr.uml/backend/utils"
	"Dr.uml/backend/utils/duerror"
	"github.com/stretchr/testify/assert"
)

func TestCreateEmptyDiagram(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		diagramType DiagramType
		expectError bool
		errorMsg    string
	}{
		{
			name:        "ValidClassDiagram",
			inputName:   "test1.uml",
			diagramType: ClassDiagram,
			expectError: false,
		},
		{
			name:        "ValidUseCaseDiagram",
			inputName:   "test2.uml",
			diagramType: UseCaseDiagram,
			expectError: true,
			errorMsg:    "Invalid diagram type",
		},
		{
			name:        "ValidSequenceDiagram",
			inputName:   "test3.uml",
			diagramType: SequenceDiagram,
			expectError: true,
			errorMsg:    "Invalid diagram type",
		},
		{
			name:        "InvalidDiagramType",
			inputName:   "test4.uml",
			diagramType: DiagramType(8),
			expectError: true,
			errorMsg:    "Invalid diagram type",
		},
		{
			name:        "InvalidDiagramType2",
			inputName:   "test5.uml",
			diagramType: DiagramType(8787),
			expectError: true,
			errorMsg:    "Invalid diagram type",
		},
		{
			name:        "InvalidName",
			inputName:   "",
			diagramType: ClassDiagram,
			expectError: true,
			errorMsg:    "file path is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := CreateEmptyUMLDiagram(tt.inputName, tt.diagramType)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
				assert.Nil(t, diagram)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, diagram)
				assert.Equal(t, tt.inputName, diagram.GetName())
				assert.Equal(t, tt.diagramType, diagram.diagramType)
				assert.WithinDuration(t, time.Now(), diagram.GetLastModified(), time.Second)
				assert.Equal(t, utils.Point{X: 0, Y: 0}, diagram.startPoint)
				// New assertions
				assert.Equal(t, "#FFFFFF", diagram.backgroundColor)
				assert.NotNil(t, diagram.componentsContainer)
			}
		})
	}
}

func TestCheckDiagramType(t *testing.T) {
	tests := []struct {
		name        string
		diagramType DiagramType
		expected    bool
	}{
		{
			name:        "ClassDiagram",
			diagramType: ClassDiagram,
			expected:    true,
		},
		{
			name:        "UseCaseDiagram",
			diagramType: UseCaseDiagram,
			expected:    false,
		},
		{
			name:        "SequenceDiagram",
			diagramType: SequenceDiagram,
			expected:    false,
		},
		{
			name:        "InvalidDiagram",
			diagramType: DiagramType(8),
			expected:    false,
		},
		{
			name:        "ZeroValue",
			diagramType: DiagramType(0),
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noError := validateDiagramType(tt.diagramType) == nil
			assert.Equal(t, tt.expected, noError)
		})
	}
}

func TestUMLDiagramGetters(t *testing.T) {
	diagram, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)
	assert.NotNil(t, diagram)

	// Test GetName
	assert.Equal(t, "TestDiagram", diagram.GetName())

	// Test GetDiagramType
	assert.Equal(t, DiagramType(ClassDiagram), diagram.GetDiagramType())

	// Test GetLastModified
	assert.WithinDuration(t, time.Now(), diagram.GetLastModified(), time.Second)

	// Test GetDrawData
	drawData := diagram.GetDrawData()
	assert.Equal(t, drawdata.Margin, drawData.Margin)
	assert.Equal(t, drawdata.LineWidth, drawData.LineWidth)
	assert.Equal(t, drawdata.DefaultDiagramColor, drawData.Color)

	err = diagram.AddGadget(component.Class, utils.Point{X: 10, Y: 20}, 0, drawdata.DefaultGadgetColor, "sample header")
	assert.NoError(t, err)
}
func TestUMLDiagram_AddGadget(t *testing.T) {
	diagram, err := CreateEmptyUMLDiagram("AddGadgetTest.uml", ClassDiagram)
	assert.NoError(t, err)
	assert.NotNil(t, diagram)
	err = diagram.AddGadget(component.Class, utils.Point{X: 5, Y: 5}, 0, drawdata.DefaultGadgetColor, "")
	assert.NoError(t, err)

	// Should have one gadget in container
	all := diagram.componentsContainer.GetAll()
	assert.Len(t, all, 1)
}

func TestUMLDiagram_AddGadget_InvalidType(t *testing.T) {
	diagram, err := CreateEmptyUMLDiagram("AddGadgetInvalid.uml", ClassDiagram)
	assert.NoError(t, err)
	assert.NotNil(t, diagram)
	// Use an invalid gadget type (assuming -1 is invalid)
	err = diagram.AddGadget(component.GadgetType(-1), utils.Point{X: 1, Y: 1}, 0, drawdata.DefaultGadgetColor, "")
	assert.Error(t, err)
}

func TestUMLDiagram_validatePoint(t *testing.T) {
	diagram, _ := CreateEmptyUMLDiagram("ValidatePoint.uml", ClassDiagram)
	assert.NoError(t, diagram.validatePoint(utils.Point{X: 0, Y: 0}))
	assert.NoError(t, diagram.validatePoint(utils.Point{X: 10, Y: 10}))
	assert.Error(t, diagram.validatePoint(utils.Point{X: -1, Y: 0}))
	assert.Error(t, diagram.validatePoint(utils.Point{X: 0, Y: -1}))
}

func TestUMLDiagram_StartAddAssociation_InvalidPoint(t *testing.T) {
	diagram, _ := CreateEmptyUMLDiagram("StartAddAssInvalid.uml", ClassDiagram)
	err := diagram.StartAddAssociation(utils.Point{X: -1, Y: 0})
	assert.Error(t, err)
}

func TestUMLDiagram_SelectAndUnselectComponent(t *testing.T) {
	diagram, _ := CreateEmptyUMLDiagram("SelectUnselect.uml", ClassDiagram)
	_ = diagram.AddGadget(component.Class, utils.Point{X: 10, Y: 10}, 0, drawdata.DefaultGadgetColor, "")

	// Select
	err := diagram.SelectComponent(utils.Point{X: 10, Y: 10})
	assert.NoError(t, err)
	assert.Len(t, diagram.componentsSelected, 1)
}

func TestUMLDiagram_UnselectAllComponents(t *testing.T) {
	diagram, _ := CreateEmptyUMLDiagram("UnselectAll.uml", ClassDiagram)
	_ = diagram.AddGadget(component.Class, utils.Point{X: 1, Y: 1}, 0, drawdata.DefaultGadgetColor, "")
	_ = diagram.SelectComponent(utils.Point{X: 1, Y: 1})
	assert.Len(t, diagram.componentsSelected, 1)
}

// func TestUMLDiagram_RegisterNotifyDrawUpdate(t *testing.T) {
// 	diagram, _ := CreateEmptyUMLDiagram("NotifyDrawUpdate.uml", ClassDiagram)
// 	err := diagram.RegisterNotifyDrawUpdate(nil)
// 	assert.Error(t, err)

// 	called := false
// 	err = diagram.RegisterNotifyDrawUpdate(func() duerror.DUError {
// 		called = true
// 		return nil
// 	})
// 	assert.NoError(t, err)
// 	_ = diagram.updateDrawData()
// 	assert.True(t, called)
// }

func TestUMLDiagram_GetDrawData(t *testing.T) {
	diagram, _ := CreateEmptyUMLDiagram("GetDrawData.uml", ClassDiagram)
	data := diagram.GetDrawData()
	assert.Equal(t, drawdata.Margin, data.Margin)
	assert.Equal(t, drawdata.LineWidth, data.LineWidth)
	assert.Equal(t, drawdata.DefaultDiagramColor, data.Color)
}

func TestUMLDiagram_RemoveSelectedComponents(t *testing.T) {
	diagram, _ := CreateEmptyUMLDiagram("RemoveSelected.uml", ClassDiagram)
	_ = diagram.AddGadget(component.Class, utils.Point{X: 2, Y: 2}, 0, drawdata.DefaultGadgetColor, "")
	_ = diagram.SelectComponent(utils.Point{X: 2, Y: 2})
	assert.Len(t, diagram.componentsSelected, 1)
	err := diagram.RemoveSelectedComponents()
	assert.NoError(t, err)
	assert.Len(t, diagram.componentsSelected, 0)
	assert.Len(t, diagram.componentsContainer.GetAll(), 0)
}

func TestUMLDiagram_EndAddAssociation_NoGadget(t *testing.T) {
	diagram, _ := CreateEmptyUMLDiagram("EndAddAssNoGadget.uml", ClassDiagram)
	_ = diagram.StartAddAssociation(utils.Point{X: 1, Y: 1})
	err := diagram.EndAddAssociation(component.AssociationType(0), utils.Point{X: 2, Y: 2})
	assert.Error(t, err)
}

func TestUMLDiagram_EndAddAssociation_InvalidPoint(t *testing.T) {
	diagram, _ := CreateEmptyUMLDiagram("EndAddAssInvalidPoint.uml", ClassDiagram)
	_ = diagram.StartAddAssociation(utils.Point{X: 1, Y: 1})
	err := diagram.EndAddAssociation(component.AssociationType(0), utils.Point{X: -1, Y: 2})
	assert.Error(t, err)
}

func TestUMLDiagram_LoadExistUMLDiagram(t *testing.T) {
	// TODO
}

func TestValidatePoint(t *testing.T) {
	diagram, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)

	// Valid point
	err = diagram.validatePoint(utils.Point{X: 10, Y: 20})
	assert.NoError(t, err)

	// Negative X
	err = diagram.validatePoint(utils.Point{X: -5, Y: 20})
	assert.Error(t, err)
	assert.Equal(t, "point coordinates must be non-negative", err.Error())

	// Negative Y
	err = diagram.validatePoint(utils.Point{X: 10, Y: -5})
	assert.Error(t, err)
	assert.Equal(t, "point coordinates must be non-negative", err.Error())

	// Both negative
	err = diagram.validatePoint(utils.Point{X: -10, Y: -10})
	assert.Error(t, err)
	assert.Equal(t, "point coordinates must be non-negative", err.Error())
}

func TestAddGadget(t *testing.T) {
	// TODO
}

func TestRemoveGadget(t *testing.T) {
	diagram, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)

	// This is a placeholder test since removeGadget is currently a no-op
	gad := &component.Gadget{}
	err = diagram.removeGadget(gad)
	assert.NoError(t, err)
}

func TestAssociationMethods(t *testing.T) {
	// TODO

	diagram, _ := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)

	gadgetPoint0 := utils.Point{X: 0, Y: 0}
	gadgetPoint1 := utils.Point{X: 0, Y: 200}
	gadgetPoint2 := utils.Point{X: 200, Y: 0}
	gadgetPoint3 := utils.Point{X: 200, Y: 200}
	diagram.AddGadget(component.Class, gadgetPoint0, 0, drawdata.DefaultGadgetColor, "")
	diagram.AddGadget(component.Class, gadgetPoint1, 0, drawdata.DefaultGadgetColor, "")
	diagram.AddGadget(component.Class, gadgetPoint2, 0, drawdata.DefaultGadgetColor, "")
	diagram.AddGadget(component.Class, gadgetPoint3, 0, drawdata.DefaultGadgetColor, "")
	gad0, _ := diagram.componentsContainer.SearchGadget(gadgetPoint0)
	gad1, _ := diagram.componentsContainer.SearchGadget(gadgetPoint1)
	gad2, _ := diagram.componentsContainer.SearchGadget(gadgetPoint2)
	gad3, _ := diagram.componentsContainer.SearchGadget(gadgetPoint3)

	t.Run("StartAddAssociation", func(t *testing.T) {
		err := diagram.StartAddAssociation(gadgetPoint0)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	var a *component.Association
	t.Run("EndAddAssociation", func(t *testing.T) {
		err := diagram.EndAddAssociation(component.Extension, gadgetPoint1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		for _, c := range diagram.componentsContainer.GetAll() {
			switch tmp := c.(type) {
			case *component.Association:
				a = tmp
			}
			if a != nil {
				break
			}
		}
		if a == nil {
			t.Errorf("add association fail")
		}
		if a.GetParentStart() != gad0 {
			t.Errorf("incorrect start parent")
		}
		if a.GetParentEnd() != gad1 {
			t.Errorf("incorrect end parent")
		}

		stList := diagram.associations[gad0][0]
		if !slices.Contains(stList, a) {
			t.Errorf("not in stList of diagram.associations")
		}
		enList := diagram.associations[gad1][1]
		if !slices.Contains(enList, a) {
			t.Errorf("not in enList of diagram.associations")
		}
	})

	t.Run("SetParentStartComponent", func(t *testing.T) {
		if a == nil {
			t.Errorf("add association fail")
		}

		dd := a.GetDrawData().(drawdata.Association)
		midPoint := utils.Point{
			X: (dd.StartX + dd.EndX) / 2,
			Y: (dd.StartY + dd.EndY) / 2,
		}
		diagram.SelectComponent(midPoint)
		c, _ := diagram.getSelectedComponent()
		if a != c.(*component.Association) {
			t.Errorf("can not select added association")
		}

		err := diagram.SetParentStartComponent(gadgetPoint2)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetParentStart() != gad2 {
			t.Errorf("incorrect start parent")
		}
		stListOld := diagram.associations[gad0][0]
		stListNew := diagram.associations[gad2][0]
		if slices.Contains(stListOld, a) {
			t.Errorf("updated association in old stList")
		}
		if !slices.Contains(stListNew, a) {
			t.Errorf("updated association not in new stList")
		}
	})

	t.Run("SetParentEndComponent", func(t *testing.T) {
		if a == nil {
			t.Errorf("add association fail")
		}

		dd := a.GetDrawData().(drawdata.Association)
		midPoint := utils.Point{
			X: (dd.StartX + dd.EndX) / 2,
			Y: (dd.StartY + dd.EndY) / 2,
		}
		diagram.selectAll(diagram.componentsSelected, false)
		diagram.SelectComponent(midPoint)
		c, _ := diagram.getSelectedComponent()
		if a != c.(*component.Association) {
			t.Errorf("can not select added association")
		}

		err := diagram.SetParentEndComponent(gadgetPoint3)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetParentEnd() != gad3 {
			t.Errorf("incorrect end parent")
		}
		enListOld := diagram.associations[gad1][1]
		enListNew := diagram.associations[gad3][1]
		if slices.Contains(enListOld, a) {
			t.Errorf("updated association in old enList")
		}
		if !slices.Contains(enListNew, a) {
			t.Errorf("updated association not in new enList")
		}
	})
}

func TestRegisterUpdateParentDraw(t *testing.T) {
	diagram, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)

	// Try to register nil function
	err = diagram.RegisterUpdateParentDraw(nil)
	assert.Error(t, err)
	assert.Equal(t, "update function cannot be nil", err.Error())

	// Register valid function
	updateCalled := false
	updateFunc := func() duerror.DUError {
		updateCalled = true
		return nil
	}

	err = diagram.RegisterUpdateParentDraw(updateFunc)
	assert.NoError(t, err)

	// Test that the function gets called through updateDrawData
	err = diagram.updateDrawData()
	assert.NoError(t, err)
	assert.True(t, updateCalled)
}

func TestUpdateDrawData(t *testing.T) {
	diagram, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)

	// Add a gadget
	err = diagram.AddGadget(component.Class, utils.Point{X: 10, Y: 20}, 0, drawdata.DefaultGadgetColor, "sample header")
	assert.NoError(t, err)

	// Check that drawData contains the gadget
	assert.Equal(t, 1, len(diagram.drawData.Gadgets))

	// Manual update of drawData
	err = diagram.updateDrawData()
	assert.NoError(t, err)

	// Test with a registered update function
	updateCalled := false
	updateFunc := func() duerror.DUError {
		updateCalled = true
		return nil
	}

	err = diagram.RegisterUpdateParentDraw(updateFunc)
	assert.NoError(t, err)

	err = diagram.updateDrawData()
	assert.NoError(t, err)
	assert.True(t, updateCalled)
}

func TestAddAttributeToGadget(t *testing.T) {
	diagram, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)

	// Try to add an attribute with no selected components
	err = diagram.AddAttributeToGadget(0, "attribute")
	assert.Error(t, err)
	assert.Equal(t, "can only operate on one component", err.Error())

	// Add a gadget to the diagram
	err = diagram.AddGadget(component.Class, utils.Point{X: 10, Y: 20}, 0, drawdata.DefaultGadgetColor, "sample header")
	assert.NoError(t, err)

	// Get the gadget from the container
	components := diagram.componentsContainer.GetAll()
	assert.Equal(t, 1, len(components))

	gadget, ok := components[0].(*component.Gadget)
	assert.True(t, ok)

	// Select the gadget
	diagram.componentsSelected[gadget] = true

	// Add attribute
	err = diagram.AddAttributeToGadget(0, "NewAttribute")
	assert.NoError(t, err)

	// Test with multiple selected components
	// First clear the selection
	diagram.selectAll(diagram.componentsSelected, false)
	assert.NoError(t, err)

	// Add a second gadget
	err = diagram.AddGadget(component.Class, utils.Point{X: 100, Y: 100}, 0, drawdata.DefaultGadgetColor, "sample header 2")
	assert.NoError(t, err)

	// Get all gadgets
	components = diagram.componentsContainer.GetAll()
	assert.Equal(t, 2, len(components))

	// Select both gadgets
	for _, comp := range components {
		diagram.componentsSelected[comp] = true
	}

	// Try to add attribute with multiple gadgets selected
	err = diagram.AddAttributeToGadget(0, "attribute")
	assert.Error(t, err)
	assert.Equal(t, "can only operate on one component", err.Error())
}

func TestLoadExistUMLDiagram(t *testing.T) {

}

func TestLoadGadgetAttributes(t *testing.T) {
	dia, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)
	assert.NotNil(t, dia)

	expectedContent := "test content"
	expectedSize := 12
	expectedStyle := attribute.Textstyle(attribute.Bold | attribute.Italic)
	expectedFontFile := os.Getenv("APP_ROOT") + "/frontend/src/assets/fonts/Inkfree.ttf"

	savedAttributeBase := utils.SavedAtt{
		Content:  expectedContent,
		Size:     expectedSize,
		Style:    int(expectedStyle),
		FontFile: expectedFontFile,
	}
	savedAttributes := make([]utils.SavedAtt, 3)
	for i := 0; i < 3; i++ {
		savedAttributes[i] = savedAttributeBase
		savedAttributes[i].Size += i
		savedAttributes[i].Ratio = 0.3 * float64(i)
	}

	gad, err := component.NewGadget(component.Class, utils.Point{}, 0, "someInvalidColorHex", "")
	assert.NoError(t, err)
	assert.NotNil(t, gad)

	err, _ = dia.loadGadgetAttributes(gad, savedAttributes) // Err-index is not important cuz we expect err is nil
	assert.NoError(t, err)

	loadedAttributes := gad.GetAttributes()

	for i := 0; i < 3; i++ {
		assert.Equal(t, 1, len(loadedAttributes[i]))
		assert.Equal(t, savedAttributes[i].Content, loadedAttributes[i][0].GetContent())
		assert.Equal(t, savedAttributes[i].Style, int(loadedAttributes[i][0].GetStyle()))
		assert.Equal(t, savedAttributes[i].FontFile, loadedAttributes[i][0].GetFontFile())
		assert.Equal(t, savedAttributes[i].Size, loadedAttributes[i][0].GetSize())
		assert.Equal(t, savedAttributes[i].FontFile, loadedAttributes[i][0].GetFontFile())
	}
}
func TestLoadGadgetAttributesButWithJsonStr(t *testing.T) {
	dia, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)
	assert.NotNil(t, dia)

	jsonStr := `{
		"Content": "test content",
		"Size": 12,
		"Style": 3,
		"FontFile": "` + os.Getenv("APP_ROOT") + `/frontend/src//assets/fonts/Inkfree.ttf"
	}`
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "/")
	var savedAttributeBase utils.SavedAtt
	err = json.Unmarshal([]byte(jsonStr), &savedAttributeBase)
	assert.NoError(t, err)

	savedAttributes := make([]utils.SavedAtt, 3)
	for i := 0; i < 3; i++ {
		savedAttributes[i] = savedAttributeBase
		savedAttributes[i].Size += i
		savedAttributes[i].Ratio = 0.3 * float64(i)
	}

	gad, err := component.NewGadget(component.Class, utils.Point{}, 0, "someInvalidColorHex", "")
	assert.NoError(t, err)
	assert.NotNil(t, gad)

	err, _ = dia.loadGadgetAttributes(gad, savedAttributes) // Err-index is not important cuz we expect err is nil
	assert.NoError(t, err)

	loadedAttributes := gad.GetAttributes()

	for i := 0; i < 3; i++ {
		assert.Equal(t, 1, len(loadedAttributes[i]))
		assert.Equal(t, savedAttributes[i].Content, loadedAttributes[i][0].GetContent())
		assert.Equal(t, savedAttributes[i].Style, int(loadedAttributes[i][0].GetStyle()))
		assert.Equal(t, savedAttributes[i].FontFile, loadedAttributes[i][0].GetFontFile())
		assert.Equal(t, savedAttributes[i].Size, loadedAttributes[i][0].GetSize())
		assert.Equal(t, savedAttributes[i].FontFile, loadedAttributes[i][0].GetFontFile())
	}
}
func TestLoadGadgets(t *testing.T) {
	dia, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)
	assert.NotNil(t, dia)

	expectedContent := "test content"
	expectedSize := 12
	expectedStyle := attribute.Textstyle(attribute.Bold | attribute.Italic)
	expectedFontFile := os.Getenv("APP_ROOT") + "/frontend/src/assets/fonts/Inkfree.ttf"

	savedAttributeBase := utils.SavedAtt{
		Content:  expectedContent,
		Size:     expectedSize,
		Style:    int(expectedStyle),
		FontFile: expectedFontFile,
	}

	savedAttributes := make([]utils.SavedAtt, 3)
	for i := 0; i < 3; i++ {
		savedAttributes[i] = savedAttributeBase
		savedAttributes[i].Size += i
		savedAttributes[i].Ratio = 0.3 * float64(i)
	}

	savedGadgetBase := utils.SavedGad{
		GadgetType: 1,
		Point:      "0, 0",
		Color:      "InvalidColorHex",
		Attributes: savedAttributes,
	}
	savedGadgets := make([]utils.SavedGad, 69)

	for i := 0; i < len(savedGadgets); i++ {
		savedGadgets[i] = savedGadgetBase
		savedGadgets[i].Layer = i
	}

	dp, err := dia.loadGadgets(savedGadgets)
	assert.NoError(t, err)
	assert.NotNil(t, dp)
	assert.Equal(t, len(savedGadgets), len(dp))
}

func TestUMLDiagram_LoadAsses(t *testing.T) {
	dia, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)
	assert.NotNil(t, dia)

	savedGadgetBase := utils.SavedGad{
		GadgetType: 1,
		Point:      "0, 0",
		Color:      "InvalidColorHex",
	}
	savedGadgets := make([]utils.SavedGad, 70)

	for i := 0; i < len(savedGadgets); i++ {
		savedGadgets[i] = savedGadgetBase
		savedGadgets[i].Layer = i
	}
	dp, err := dia.loadGadgets(savedGadgets)
	assert.NoError(t, err)

	expectedAssType := component.AssociationType(1)
	expectedLayer := 0
	expectedStartRatio := [2]float64{0.1, 0.2}
	expectedEndRatio := [2]float64{0.3, 0.4}

	savedAssBase := utils.SavedAss{
		AssType:         int(expectedAssType),
		Layer:           expectedLayer,
		StartPointRatio: expectedStartRatio,
		EndPointRatio:   expectedEndRatio,
	}
	savedAsses := make([]utils.SavedAss, 69)
	for i := 0; i < len(savedAsses); i++ {
		savedAsses[i] = savedAssBase
		savedAsses[i].Parents = []int{i, i + 1} // Assuming each association connects
	}
	err = dia.loadAsses(savedAsses, dp)
	assert.NoError(t, err)
	// Check if associations are loaded correctly
	components := dia.componentsContainer.GetAll()
	for _, comp := range components {
		switch comp.(type) {
		case *component.Association:
			assert.Equal(t, expectedAssType, comp.(*component.Association).GetAssType())
			assert.Equal(t, expectedLayer, comp.(*component.Association).GetLayer())
		default:
			continue
		}
	}
}

func TestUMLDiagram_loadAssAttributes(t *testing.T) {
	dia, err := CreateEmptyUMLDiagram("TestDiagram", ClassDiagram)
	assert.NoError(t, err)
	assert.NotNil(t, dia)

	var parents = [2]*component.Gadget{nil, nil}
	for i := 0; i < 2; i++ {
		parents[i], err = component.NewGadget(component.Class, utils.Point{X: 0, Y: 0}, 0, "InvalidColorHex", "")
		assert.NoError(t, err)
	}

	ass, err := component.NewAssociation(parents, component.AssociationType(1), utils.Point{X: 0, Y: 0}, utils.Point{X: 1, Y: 1})
	assert.NoError(t, err)
	assert.NotNil(t, ass)

	// Prepare attributes

	expectedContent := "test"
	expectedStyle := int(attribute.Bold)
	expectedFontFile := os.Getenv("APP_ROOT") + "/frontend/src/assets/fonts/Inkfree.ttf"
	expectedRatio := 0.69

	savedAttBase := utils.SavedAtt{
		Content:  expectedContent,
		Style:    expectedStyle,
		FontFile: expectedFontFile,
		Ratio:    expectedRatio,
	}

	attributes := make([]utils.SavedAtt, 2)
	for i := 0; i < len(attributes); i++ {
		attributes[i] = savedAttBase
		attributes[i].Size = i + 1
	}

	// Should succeed
	errRet, idx := dia.loadAssAttributes(ass, attributes)
	assert.NoError(t, errRet)
	assert.Equal(t, 0, idx)

	atts := ass.GetAttributes()
	for i, att := range atts {
		assert.Equal(t, expectedContent, att.GetContent())
		assert.Equal(t, expectedStyle, int(att.GetStyle()))
		assert.Equal(t, expectedRatio, att.GetRatio())
		assert.Equal(t, expectedFontFile, att.GetFontFile())
		assert.Equal(t, i+1, att.GetSize())
	}

	// Should fail with nil association
	errRet, idx = dia.loadAssAttributes(nil, attributes)
	assert.Error(t, errRet)
	assert.Equal(t, 0, idx)
}

func TestUMLDiagram_SaveToFile(t *testing.T) {
	diagram, err := CreateEmptyUMLDiagram("SaveToFileTest.uml", ClassDiagram)
	assert.NoError(t, err)
	assert.NotNil(t, diagram)

	// Add two gadgets
	err = diagram.AddGadget(component.Class, utils.Point{X: 1, Y: 2}, 0, drawdata.DefaultGadgetColor, "Header1")
	assert.NoError(t, err)
	err = diagram.AddGadget(component.Class, utils.Point{X: 3, Y: 4}, 1, drawdata.DefaultGadgetColor, "Header2")
	assert.NoError(t, err)

	// Add an association between the two gadgets
	gadgets := diagram.componentsContainer.GetAll()
	assert.Len(t, gadgets, 2)
	gad1, ok1 := gadgets[0].(*component.Gadget)
	gad2, ok2 := gadgets[1].(*component.Gadget)
	assert.True(t, ok1)
	assert.True(t, ok2)
	ass, err := component.NewAssociation([2]*component.Gadget{gad1, gad2}, component.AssociationType(1), utils.Point{X: 1, Y: 2}, utils.Point{X: 3, Y: 4})
	assert.NoError(t, err)
	assert.NotNil(t, ass)
	err = diagram.componentsContainer.Insert(ass)
	assert.NoError(t, err)
	diagram.associations[gad1] = [2][]*component.Association{{ass}, {}}
	diagram.associations[gad2] = [2][]*component.Association{{}, {ass}}

	// Save to file
	_, err = diagram.SaveToFile("SaveToFileTest.uml")
	assert.NoError(t, err)
}

// Mock container for testing selection methods
type mockContainer struct {
	mockComponent component.Component
	pointToMatch  utils.Point
}

func (m *mockContainer) Insert(c component.Component) duerror.DUError {
	return nil
}

func (m *mockContainer) Remove(c component.Component) duerror.DUError {
	return nil
}

func (m *mockContainer) Search(p utils.Point) (component.Component, duerror.DUError) {
	if p.X == m.pointToMatch.X && p.Y == m.pointToMatch.Y {
		return m.mockComponent, nil
	}
	return nil, nil
}

func (m *mockContainer) GetAll() []component.Component {
	if m.mockComponent != nil {
		return []component.Component{m.mockComponent}
	}
	return []component.Component{}
}

func (m *mockContainer) Len() (int, duerror.DUError) {
	if m.mockComponent != nil {
		return 1, nil
	}
	return 0, nil
}

func CMD_ADD_GADGET(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)
	header := "test gadget"
	d.AddGadget(component.Class, utils.Point{X: 0, Y: 0}, 0, drawdata.DefaultGadgetColor, header)

	t.Run("undo add gadget", func(t *testing.T) {
		err := d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if d.componentsContainer.Len() != 0 {
			t.Errorf("fail to remove gadget from component container")
		}
	})

	t.Run("redo add gadget", func(t *testing.T) {
		err := d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if d.componentsContainer.Len() != 1 {
			t.Errorf("fail to recover component container")
		}
		dd := d.GetDrawData()
		content := dd.Gadgets[0].Attributes[0][0].Content
		if content != header {
			t.Errorf("fail to recover gadget content: %v, got %v", header, content)
		}
	})
}

func CMD_ADD_ASSOCIATION(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	gadPoint1 := utils.Point{X: 200, Y: 200}
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, "test gadget0")
	d.AddGadget(component.Class, gadPoint1, 0, drawdata.DefaultGadgetColor, "test gadget1")

	assType := component.Composition
	d.StartAddAssociation(gadPoint0)
	d.EndAddAssociation(component.AssociationType(assType), gadPoint1)

	t.Run("undo add association", func(t *testing.T) {
		err := d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if d.componentsContainer.Len() != 2 {
			t.Errorf("fail to remove association from component container")
		}
		dd := d.GetDrawData()
		assLen := len(dd.Associations)
		if assLen != 0 {
			t.Errorf("fail to remove correct component")
		}
	})

	t.Run("redo add association", func(t *testing.T) {
		err := d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if d.componentsContainer.Len() != 3 {
			t.Errorf("fail to recover component container")
		}
		dd := d.GetDrawData()
		assLen := len(dd.Associations)
		if assLen != 1 {
			t.Errorf("fail to recover association")
		}
		newAssType := dd.Associations[0].AssType
		if newAssType != assType {
			t.Errorf("fail to recover association type: %v, got %v", assType, newAssType)
		}
	})
}

func CMD_REMOVE_SELECTED_COMPONENTS(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	gadPoint1 := utils.Point{X: 200, Y: 200}
	gadPoint2 := utils.Point{X: 400, Y: 400}
	header0 := "test gadget0"
	header1 := "test gadget1"
	header2 := "test gadget2"
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, header0)
	d.AddGadget(component.Class, gadPoint1, 0, drawdata.DefaultGadgetColor, header1)
	d.AddGadget(component.Class, gadPoint2, 0, drawdata.DefaultGadgetColor, header2)

	assType0 := component.Composition
	assType1 := component.Dependency
	assType2 := component.Extension
	d.StartAddAssociation(gadPoint0)
	d.EndAddAssociation(component.AssociationType(assType0), gadPoint1)
	d.StartAddAssociation(gadPoint0)
	d.EndAddAssociation(component.AssociationType(assType1), gadPoint0)
	d.StartAddAssociation(gadPoint2)
	d.EndAddAssociation(component.AssociationType(assType2), gadPoint2)

	t.Run("select and remove components", func(t *testing.T) {
		d.SelectComponent(gadPoint0)
		d.SelectComponent(gadPoint1)
		midPoint := utils.Point{
			X: (gadPoint0.X + gadPoint1.X) / 2,
			Y: (gadPoint0.Y + gadPoint1.Y) / 2,
		}
		err := d.SelectComponent(midPoint)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 3 {
			t.Errorf("fail to select components")
		}

		err = d.RemoveSelectedComponents()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		dd := d.GetDrawData()
		if len(dd.Gadgets) != 1 {
			t.Errorf("failed to remove gadget")
		}
		content := dd.Gadgets[0].Attributes[0][0].Content
		if content != header2 {
			t.Errorf("failed to remove gadget")
		}

		if len(dd.Associations) != 1 {
			t.Errorf("failed to remove association")
		}
		assType := dd.Associations[0].AssType
		if assType != assType2 {
			t.Errorf("failed to remove association")
		}
	})

	t.Run("undo remove select components", func(t *testing.T) {
		err := d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(d.componentsSelected) != 3 {
			t.Errorf("fail to recover componentsSelected")
		}

		dd := d.GetDrawData()
		if len(dd.Gadgets) != 3 {
			t.Errorf("failed to recover gadget")
		} else {
			headers := map[string]bool{header0: true, header1: true, header2: true}
			for i := 0; i < 3; i++ {
				content := dd.Gadgets[i].Attributes[0][0].Content
				val, ok := headers[content]
				if !ok || !val {
					t.Errorf("failed to recover gadget. incorrect header: %v", content)
				}
				headers[content] = false
			}
		}

		if len(dd.Associations) != 3 {
			t.Errorf("failed to recover association")
		} else {
			types := map[int]bool{assType0: true, assType1: true, assType2: true}
			for i := 0; i < 3; i++ {
				assType := dd.Associations[i].AssType
				val, ok := types[assType]
				if !ok || !val {
					t.Errorf("failed to recover association. incorrect type: %v", assType)
				}
				types[assType] = false
			}
		}
	})

	t.Run("redo remove select components", func(t *testing.T) {
		err := d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		dd := d.GetDrawData()
		if len(dd.Gadgets) != 1 {
			t.Errorf("failed to remove gadget")
		}
		content := dd.Gadgets[0].Attributes[0][0].Content
		if content != header2 {
			t.Errorf("failed to remove gadget")
		}

		if len(dd.Associations) != 1 {
			t.Errorf("failed to remove association")
		}
		assType := dd.Associations[0].AssType
		if assType != assType2 {
			t.Errorf("failed to remove association")
		}
	})
}

func CMD_SELECT_COMPONENT(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	gadPoint1 := utils.Point{X: 200, Y: 200}
	header0 := "test gadget0"
	header1 := "test gadget1"
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, header0)
	d.AddGadget(component.Class, gadPoint1, 0, drawdata.DefaultGadgetColor, header1)

	assType0 := component.Composition
	d.StartAddAssociation(gadPoint0)
	d.EndAddAssociation(component.AssociationType(assType0), gadPoint1)

	midPoint := utils.Point{
		X: (gadPoint0.X + gadPoint1.X) / 2,
		Y: (gadPoint0.Y + gadPoint1.Y) / 2,
	}
	noWherePoint := utils.Point{X: 1000, Y: 1000}

	t.Run("select components", func(t *testing.T) {
		d.SelectComponent(gadPoint0)
		if len(d.componentsSelected) != 1 {
			t.Errorf("select one fails")
		}
		d.SelectComponent(gadPoint1)
		if len(d.componentsSelected) != 2 {
			t.Errorf("select two fails")
		}
		d.SelectComponent(midPoint)
		if len(d.componentsSelected) != 3 {
			t.Errorf("select three fails")
		}
		d.SelectComponent(gadPoint0)
		if len(d.componentsSelected) != 3 {
			t.Errorf("select already selected association fails")
		}
		d.SelectComponent(noWherePoint)
		if len(d.componentsSelected) != 0 {
			t.Errorf("select empty place fails")
		}
	})

	t.Run("undo select components", func(t *testing.T) {
		err := d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 3 {
			t.Errorf("undo select empty place fails")
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 2 {
			t.Errorf("undo select three fails")
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 1 {
			t.Errorf("undo select two fails")
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 0 {
			t.Errorf("undo select one fails")
		}
	})

	t.Run("redo select components", func(t *testing.T) {
		err := d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 1 {
			t.Errorf("select one fails")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 2 {
			t.Errorf("select two fails")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 3 {
			t.Errorf("select three fails")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 3 {
			t.Errorf("select already selected association fails")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(d.componentsSelected) != 0 {
			t.Errorf("select empty place fails")
		}
	})
}

func CMD_COMPONENT_SETTER(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	gadPoint1 := utils.Point{X: 200, Y: 200}
	header0 := "test gadget0"
	header1 := "test gadget1"
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, header0)
	d.AddGadget(component.Class, gadPoint1, 0, drawdata.DefaultGadgetColor, header1)

	assType0 := component.Composition
	d.StartAddAssociation(gadPoint0)
	d.EndAddAssociation(component.AssociationType(assType0), gadPoint1)
	midPoint := utils.Point{
		X: (gadPoint0.X + gadPoint1.X) / 2,
		Y: (gadPoint0.Y + gadPoint1.Y) / 2,
	}

	t.Run("setter gadget", func(t *testing.T) {
		newLayer := 69
		newColor := "#123456"
		d.SelectComponent(gadPoint0)
		err := d.SetLayerComponent(newLayer)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		err = d.SetColorComponent(newColor)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, gdd := range d.GetDrawData().Gadgets {
			if gdd.Attributes[0][0].Content == header0 {
				if gdd.Layer != newLayer {
					t.Errorf("unexpected layer: %v, got %v", newLayer, gdd.Layer)
				}
				if gdd.Color != newColor {
					t.Errorf("unexpected color: %v, got %v", newColor, gdd.Color)
				}
			}
		}

		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, gdd := range d.GetDrawData().Gadgets {
			if gdd.Attributes[0][0].Content == header0 {
				if gdd.Layer == newLayer {
					t.Errorf("undo set layer fails")
				}
				if gdd.Color == newColor {
					t.Errorf("undo set color fails")
				}
			}
		}

		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, gdd := range d.GetDrawData().Gadgets {
			if gdd.Attributes[0][0].Content == header0 {
				if gdd.Layer != newLayer {
					t.Errorf("redo set layer fails")
				}
				if gdd.Color != newColor {
					t.Errorf("redo set color fails")
				}
			}
		}
	})

	t.Run("setter gadget", func(t *testing.T) {
		newLayer := 69
		d.SelectComponent(utils.Point{X: 1000, Y: 1000})
		d.SelectComponent(midPoint)
		err := d.SetLayerComponent(newLayer)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, add := range d.GetDrawData().Associations {
			if add.Layer != newLayer {
				t.Errorf("set association layer fails")
			}
		}

		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, add := range d.GetDrawData().Associations {
			if add.Layer == newLayer {
				t.Errorf("undo set association layer fails")
			}
		}

		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, add := range d.GetDrawData().Associations {
			if add.Layer != newLayer {
				t.Errorf("redo set association layer fails")
			}
		}
	})
}

func CMD_MOVE_GADGET(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	gadPoint1 := utils.Point{X: 200, Y: 200}
	header0 := "test gadget0"
	header1 := "test gadget1"
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, header0)
	d.AddGadget(component.Class, gadPoint1, 0, drawdata.DefaultGadgetColor, header1)

	assType0 := component.Composition
	d.StartAddAssociation(gadPoint0)
	d.EndAddAssociation(component.AssociationType(assType0), gadPoint1)

	add := d.GetDrawData().Associations[0]
	oldAssStart := utils.Point{X: add.StartX, Y: add.StartY}

	newPoint := utils.Point{X: 69, Y: 69}
	t.Run("move gadget", func(t *testing.T) {
		d.SelectComponent(gadPoint0)
		err := d.SetPointComponent(newPoint)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, gdd := range d.GetDrawData().Gadgets {
			if gdd.Attributes[0][0].Content == header0 {
				point := utils.Point{X: gdd.X, Y: gdd.Y}
				if point != newPoint {
					t.Errorf("set point fails")
				}
			}
		}
		add := d.GetDrawData().Associations[0]
		point := utils.Point{X: add.StartX, Y: add.StartY}
		if point == oldAssStart {
			t.Errorf("move gadget doesnt change its association")
		}
	})

	t.Run("undo move gadget", func(t *testing.T) {
		err := d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, gdd := range d.GetDrawData().Gadgets {
			if gdd.Attributes[0][0].Content == header0 {
				point := utils.Point{X: gdd.X, Y: gdd.Y}
				if point == newPoint {
					t.Errorf("undo set point fails")
				}
			}
		}
		add := d.GetDrawData().Associations[0]
		point := utils.Point{X: add.StartX, Y: add.StartY}
		if point != oldAssStart {
			t.Errorf("move gadget doesnt change its association")
		}
	})

	t.Run("redo move gadget", func(t *testing.T) {
		err := d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, gdd := range d.GetDrawData().Gadgets {
			if gdd.Attributes[0][0].Content == header0 {
				point := utils.Point{X: gdd.X, Y: gdd.Y}
				if point != newPoint {
					t.Errorf("redo set point fails")
				}
			}
		}
		add := d.GetDrawData().Associations[0]
		point := utils.Point{X: add.StartX, Y: add.StartY}
		if point == oldAssStart {
			t.Errorf("move gadget doesnt change its association")
		}
	})
}

func CMD_SET_PARENT(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	gadPoint1 := utils.Point{X: 200, Y: 200}
	gadPoint2 := utils.Point{X: 400, Y: 400}
	header0 := "test gadget0"
	header1 := "test gadget1"
	header2 := "test gadget2"
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, header0)
	d.AddGadget(component.Class, gadPoint1, 0, drawdata.DefaultGadgetColor, header1)
	d.AddGadget(component.Class, gadPoint2, 0, drawdata.DefaultGadgetColor, header2)

	assType0 := component.Composition
	d.StartAddAssociation(gadPoint0)
	d.EndAddAssociation(component.AssociationType(assType0), gadPoint1)

	midPoint := utils.Point{
		X: (gadPoint0.X + gadPoint1.X) / 2,
		Y: (gadPoint0.Y + gadPoint1.Y) / 2,
	}
	d.SelectComponent(midPoint)
	c, _ := d.componentsContainer.Search(midPoint)
	a := c.(*component.Association)

	g0, _ := d.componentsContainer.SearchGadget(gadPoint0)
	g1, _ := d.componentsContainer.SearchGadget(gadPoint1)
	g2, _ := d.componentsContainer.SearchGadget(gadPoint2)
	t.Run("set parent start", func(t *testing.T) {
		err := d.SetParentStartComponent(gadPoint2)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetParentStart() != g2 {
			t.Errorf("set parent start fails")
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetParentStart() != g0 {
			t.Errorf("undo set parent start fails")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetParentStart() != g2 {
			t.Errorf("undi set parent start fails")
		}
	})

	t.Run("set parent end", func(t *testing.T) {
		err := d.SetParentEndComponent(gadPoint2)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetParentStart() != g2 {
			t.Errorf("set parent start fails")
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetParentStart() != g1 {
			t.Errorf("undo set parent start fails")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetParentStart() != g2 {
			t.Errorf("undi set parent start fails")
		}
	})
}

func CMD_SET_ATTR_GAD(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	header0 := "test gadget0"
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, header0)
	d.SelectComponent(gadPoint0)

	g, _ := d.componentsContainer.SearchGadget(gadPoint0)

	t.Run("set content", func(t *testing.T) {
		getValue := func(g *component.Gadget) string {
			return g.GetDrawData().(drawdata.Gadget).Attributes[0][0].Content
		}
		newValue := "new content"
		oldValue := getValue(g)

		err := d.SetAttrContentComponent(0, 0, newValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(g) != newValue {
			t.Errorf("unexpected value: %v, got %v", newValue, getValue(g))
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(g) != oldValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(g))
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(g) != newValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(g))
		}
	})

	t.Run("set size", func(t *testing.T) {
		getValue := func(g *component.Gadget) int {
			return g.GetDrawData().(drawdata.Gadget).Attributes[0][0].FontSize
		}
		newValue := 69
		oldValue := getValue(g)

		err := d.SetAttrSizeComponent(0, 0, newValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(g) != newValue {
			t.Errorf("unexpected value: %v, got %v", newValue, getValue(g))
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(g) != oldValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(g))
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(g) != newValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(g))
		}
	})

	t.Run("set style", func(t *testing.T) {
		getValue := func(g *component.Gadget) int {
			return g.GetDrawData().(drawdata.Gadget).Attributes[0][0].FontStyle
		}
		newValue := attribute.Bold
		oldValue := getValue(g)

		err := d.SetAttrStyleComponent(0, 0, newValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(g) != newValue {
			t.Errorf("unexpected value: %v, got %v", newValue, getValue(g))
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(g) != oldValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(g))
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(g) != newValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(g))
		}
	})
}

func CMD_SET_ATTR_ASS(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	gadPoint1 := utils.Point{X: 200, Y: 200}
	header0 := "test gadget0"
	header1 := "test gadget1"
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, header0)
	d.AddGadget(component.Class, gadPoint1, 0, drawdata.DefaultGadgetColor, header1)

	assType0 := component.Composition
	d.StartAddAssociation(gadPoint0)
	d.EndAddAssociation(component.AssociationType(assType0), gadPoint1)

	midPoint := utils.Point{
		X: (gadPoint0.X + gadPoint1.X) / 2,
		Y: (gadPoint0.Y + gadPoint1.Y) / 2,
	}
	d.SelectComponent(midPoint)
	c, _ := d.componentsContainer.Search(midPoint)
	a := c.(*component.Association)

	d.AddAttributeToAssociation(0.5, "test")

	t.Run("set content", func(t *testing.T) {
		getValue := func(a *component.Association) string {
			return a.GetDrawData().(drawdata.Association).Attributes[0].Content
		}
		newValue := "new content"
		oldValue := getValue(a)

		err := d.SetAttrContentComponent(0, 0, newValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != newValue {
			t.Errorf("unexpected value: %v, got %v", newValue, getValue(a))
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != oldValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(a))
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != newValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(a))
		}
	})

	t.Run("set size", func(t *testing.T) {
		getValue := func(a *component.Association) int {
			return a.GetDrawData().(drawdata.Association).Attributes[0].FontSize
		}
		newValue := 69
		oldValue := getValue(a)

		err := d.SetAttrSizeComponent(0, 0, newValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != newValue {
			t.Errorf("unexpected value: %v, got %v", newValue, getValue(a))
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != oldValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(a))
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != newValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(a))
		}
	})

	t.Run("set style", func(t *testing.T) {
		getValue := func(a *component.Association) int {
			return a.GetDrawData().(drawdata.Association).Attributes[0].FontStyle
		}
		newValue := attribute.Bold
		oldValue := getValue(a)

		err := d.SetAttrStyleComponent(0, 0, newValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != newValue {
			t.Errorf("unexpected value: %v, got %v", newValue, getValue(a))
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != oldValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(a))
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != newValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(a))
		}
	})

	t.Run("set ratio", func(t *testing.T) {
		getValue := func(a *component.Association) float64 {
			return a.GetDrawData().(drawdata.Association).Attributes[0].Ratio
		}
		newValue := 0.69
		oldValue := getValue(a)

		err := d.SetAttrRatioComponent(0, 0, newValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != newValue {
			t.Errorf("unexpected value: %v, got %v", newValue, getValue(a))
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != oldValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(a))
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if getValue(a) != newValue {
			t.Errorf("unexpected value: %v, got %v", oldValue, getValue(a))
		}
	})
}

func CMD_ADD_REMOVE_ATTR_GAD(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	header0 := "test gadget0"
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, header0)
	d.SelectComponent(gadPoint0)

	g, _ := d.componentsContainer.SearchGadget(gadPoint0)
	section := 1
	t.Run("add attribute", func(t *testing.T) {
		err := d.AddAttributeToGadget(section, "ඞ")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if g.GetAttributesLen()[section] != 1 {
			t.Errorf("add attr to gad fail")
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if g.GetAttributesLen()[section] != 0 {
			t.Errorf("undo add attr to gad fail")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if g.GetAttributesLen()[section] != 1 {
			t.Errorf("undo add attr to gad fail")
		}
	})

	t.Run("remove attribute", func(t *testing.T) {
		err := d.RemoveAttributeFromGadget(section, 0)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if g.GetAttributesLen()[section] != 0 {
			t.Errorf("remove attr to gad fail")
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if g.GetAttributesLen()[section] != 1 {
			t.Errorf("undo remove attr to gad fail")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if g.GetAttributesLen()[section] != 0 {
			t.Errorf("undo remove attr to gad fail")
		}
	})
}

func CMD_ADD_REMOVE_ATTR_ASS(t *testing.T) {
	d, _ := CreateEmptyUMLDiagram("test.uml", ClassDiagram)

	gadPoint0 := utils.Point{X: 0, Y: 0}
	gadPoint1 := utils.Point{X: 200, Y: 200}
	header0 := "test gadget0"
	header1 := "test gadget1"
	d.AddGadget(component.Class, gadPoint0, 0, drawdata.DefaultGadgetColor, header0)
	d.AddGadget(component.Class, gadPoint1, 0, drawdata.DefaultGadgetColor, header1)

	assType0 := component.Composition
	d.StartAddAssociation(gadPoint0)
	d.EndAddAssociation(component.AssociationType(assType0), gadPoint1)

	midPoint := utils.Point{
		X: (gadPoint0.X + gadPoint1.X) / 2,
		Y: (gadPoint0.Y + gadPoint1.Y) / 2,
	}
	d.SelectComponent(midPoint)
	c, _ := d.componentsContainer.Search(midPoint)
	a := c.(*component.Association)

	t.Run("add attribute", func(t *testing.T) {
		err := d.AddAttributeToAssociation(0.5, "ඞ")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetAttributesLen() != 1 {
			t.Errorf("add attr to ass fail")
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetAttributesLen() != 0 {
			t.Errorf("undo add attr to ass fail")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetAttributesLen() != 1 {
			t.Errorf("undo add attr to ass fail")
		}
	})

	t.Run("remove attribute", func(t *testing.T) {
		err := d.RemoveAttributeFromAssociation(0)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetAttributesLen() != 0 {
			t.Errorf("remove attr to ass fail")
		}
		err = d.Undo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetAttributesLen() != 1 {
			t.Errorf("undo remove attr to ass fail")
		}
		err = d.Redo()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if a.GetAttributesLen() != 0 {
			t.Errorf("undo remove attr to ass fail")
		}
	})
}
