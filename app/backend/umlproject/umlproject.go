package umlproject

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"time"

	"Dr.uml/backend/component"
	"Dr.uml/backend/drawdata"
	"Dr.uml/backend/umldiagram"
	"Dr.uml/backend/utils"
	"Dr.uml/backend/utils/duerror"

	"github.com/titanous/json5"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type UMLProject struct {
	ctx               context.Context
	name              string
	lastModified      time.Time
	currentDiagram    *umldiagram.UMLDiagram            // The currently selected diagram
	availableDiagrams map[string]bool                   // Use a map to store diagrams, keyed by their ID
	activeDiagrams    map[string]*umldiagram.UMLDiagram // Keep track of active diagrams
	runFrontend       bool
}

// Constructor
func CreateEmptyUMLProject(fileName string) (*UMLProject, duerror.DUError) {
	// TODO: also check the file is exist or not
	if err := utils.ValidateFilePath(fileName); err != nil {
		return nil, err
	}
	return &UMLProject{
		name:              fileName,
		lastModified:      time.Now(),
		availableDiagrams: make(map[string]bool),
		activeDiagrams:    make(map[string]*umldiagram.UMLDiagram),
	}, nil
}

func LoadExistUMLProject(fileName string) (*UMLProject, duerror.DUError) {
	// TODO
	return nil, nil
}

// Other functions
func (p *UMLProject) Startup(ctx context.Context) {
	p.ctx = ctx
	// should not AddNewDiagram and SelectDiagram here
	// TODO: Remove this
	p.runFrontend = true
	p.CreateEmptyUMLDiagram(umldiagram.ClassDiagram, "new class diagram")
	p.SelectDiagram("new class diagram")
}

// Getter
func (p *UMLProject) GetName() string {
	return p.name
}

func (p *UMLProject) GetLastModified() time.Time {
	return p.lastModified
}

func (p *UMLProject) GetCurrentDiagramName() string {
	if p.currentDiagram == nil {
		return ""
	}
	return p.currentDiagram.GetName()
}

func (p *UMLProject) GetAvailableDiagramsNames() []string {
	return slices.Collect(maps.Keys(p.availableDiagrams))
}

func (p *UMLProject) GetActiveDiagramsNames() []string {
	activeNames := make([]string, 0, len(p.activeDiagrams))
	for _, d := range p.activeDiagrams {
		activeNames = append(activeNames, d.GetName())
	}
	return activeNames
}

// Other methods
// SelectDiagram makes a diagram to be the current diagram, loads it if not loaded yet
func (p *UMLProject) SelectDiagram(diagramName string) duerror.DUError {
	if _, ok := p.availableDiagrams[diagramName]; !ok {
		return duerror.NewInvalidArgumentError("Diagram not found")
	}
	if _, ok := p.activeDiagrams[diagramName]; !ok {
		err := p.OpenDiagram(diagramName)
		if err != nil {
			return err
		}
	}
	p.currentDiagram = p.activeDiagrams[diagramName]
	return nil
}

func (p *UMLProject) CreateEmptyUMLDiagram(diagramType umldiagram.DiagramType, diagramName string) duerror.DUError {
	if _, ok := p.availableDiagrams[diagramName]; ok {
		return duerror.NewInvalidArgumentError("Diagram name already exists")
	}
	d, err := umldiagram.CreateEmptyUMLDiagram(diagramName, diagramType)
	if err != nil {
		return err
	}
	p.availableDiagrams[diagramName] = true
	p.activeDiagrams[diagramName] = d
	p.lastModified = time.Now()
	return nil
}

func (p *UMLProject) CloseDiagram(diagramName string) duerror.DUError {
	// TODO: save file?
	if _, ok := p.activeDiagrams[diagramName]; !ok {
		return duerror.NewInvalidArgumentError("Diagram not loaded")
	}
	if p.currentDiagram != nil && p.currentDiagram.GetName() == diagramName {
		p.currentDiagram = nil
	}
	delete(p.activeDiagrams, diagramName)
	return nil
}

func (p *UMLProject) DeleteDiagram(diagramName string) duerror.DUError {
	// TODO: remove the file
	return nil
}

func (p *UMLProject) AddGadget(gadgetType component.GadgetType, point utils.Point) duerror.DUError {
	if p.currentDiagram == nil {
		return duerror.NewInvalidArgumentError("No current diagram selected")
	}
	if err := p.currentDiagram.AddGadget(gadgetType, point); err != nil {
		return err
	}
	p.lastModified = time.Now()
	return p.InvalidateCanvas()
}

func (p *UMLProject) StartAddAssociation(point utils.Point) duerror.DUError {
	if p.currentDiagram == nil {
		return duerror.NewInvalidArgumentError("No current diagram selected")
	}
	return p.currentDiagram.StartAddAssociation(point)
}

func (p *UMLProject) EndAddAssociation(associationType component.AssociationType, point utils.Point) duerror.DUError {
	if p.currentDiagram == nil {
		return duerror.NewInvalidArgumentError("No current diagram selected")
	}
	if err := p.currentDiagram.EndAddAssociation(associationType, point); err != nil {
		return err
	}
	p.lastModified = time.Now()
	return p.InvalidateCanvas()
}

func (p *UMLProject) RemoveSelectedComponents() duerror.DUError {
	if p.currentDiagram == nil {
		return duerror.NewInvalidArgumentError("No current diagram selected")
	}
	if err := p.currentDiagram.RemoveSelectedComponents(); err != nil {
		return err
	}
	p.lastModified = time.Now()
	return nil
}

// Draw
func (p *UMLProject) GetDrawData() drawdata.Diagram {
	if p.currentDiagram == nil {
		return drawdata.Diagram{}
	}
	return p.currentDiagram.GetDrawData()
}

func (p *UMLProject) InvalidateCanvas() duerror.DUError {
	if !p.runFrontend {
		return nil
	}

	if p.currentDiagram == nil {
		return duerror.NewInvalidArgumentError("No current diagram selected")
	}
	// p.notifyDrawUpdate(p.currentDiagram.GetName())
	// log.Println("InvalidateCanvas")
	runtime.EventsEmit(p.ctx, "backend-event", p.GetDrawData())
	return nil
}

func (p *UMLProject) OpenDiagram(filename string) duerror.DUError {
	err := utils.ValidateFilePath(filename)
	if err != nil {
		return err
	}

	if p.availableDiagrams[filename] {
		return p.SelectDiagram(filename)
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return duerror.NewFileIOError(fmt.Sprintf("Failed to open file %s.\n Error: %s", filename, err.Error()))
	}
	defer file.Close()

	decoder := json5.NewDecoder(file)
	var savedFileData utils.SavedFile
	{
	}
	if err := decoder.Decode(&savedFileData); err != nil {
		return duerror.NewInvalidArgumentError(fmt.Sprintf("Failed to decode file %s.\n Error: %s", filename, err.Error()))
	}
	switch savedFileData.Filetype {
	case utils.FiletypeDiagram:
		dia, err := umldiagram.LoadExistUMLDiagram(savedFileData)
		if err != nil {
			return err
		}
		p.availableDiagrams
		break
	case utils.FiletypeSubmodule:

		// TODO
		break
	default:
		return duerror.NewInvalidArgumentError(fmt.Sprintf("Unknown filetype %d", savedFileData.Filetype))
	}

	return nil
}
