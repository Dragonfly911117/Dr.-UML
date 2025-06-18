package utils

import (
	"Dr.uml/backend/drawdata"
	"errors"
	"github.com/labstack/gommon/log"
	"image"
	"image/color"
	"image/png"
	"os"

	"Dr.uml/backend/utils/duerror"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func loadFont(file string) (*opentype.Font, duerror.DUError) {
	fontBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, duerror.NewFileIOError(err.Error())
	}
	fnt, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, duerror.NewFileIOError(err.Error())
	}
	return fnt, nil
}

// genGlyphImage is a helper function for the ease of development and should not
// be used on production.
func genGlyphImage(face font.Face, str string) (err duerror.DUError) { // notest
	// Draw the string to an image
	imgWidth := 400
	imgHeight := 100
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  fixed.P(10, 50),
	}
	d.DrawString(str)

	outFile, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer func() {
		err = errors.Join(outFile.Close())
	}()

	if err := png.Encode(outFile, img); err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}

	return err
}

func GetTextSize(str string, size int, fontFile string) (height int, width int, err duerror.DUError) {
	defaultFontFile := os.Getenv("APP_ROOT") + drawdata.DefaultAttributeFontFile
	if fontFile == "" {
		fontFile = defaultFontFile
	}
	dpi := 72
	if size <= 0 {
		return 0, 0, duerror.NewInvalidArgumentError("size must be greater than 0")
	}
	fnt, err := loadFont(fontFile)
	if err != nil {
		return 0, 0, err
	}
	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     float64(dpi),
		Hinting: font.HintingFull,
	})
	if err != nil { //notest
		return 0, 0, duerror.NewFileIOError(err.Error())
	}
	defer func() {
		err = errors.Join(face.Close())
	}()

	//genGlyphImage(face, str)

	var fixedWidth fixed.Int26_6
	for _, r := range str {
		advance, ok := face.GlyphAdvance(r)
		if !ok { // notest
			return 0, 0, duerror.NewFileIOError("glyph not found")
		}
		fixedWidth += advance
	}
	metrics := face.Metrics()
	fixedHeight := metrics.Ascent + metrics.Descent
	return fixedHeight.Round(), fixedWidth.Round(), nil
}
