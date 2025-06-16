package utils

import (
	"Dr.uml/backend/drawdata"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GetTextSize(t *testing.T) {
	tests := []struct {
		name           string
		str            string
		size           int
		fontFile       string
		expectedHeight int
		expectedWidth  int
		hasError       bool
	}{
		{
			name:           "correct",
			str:            "Hello, World!",
			size:           12,
			fontFile:       os.Getenv("APP_ROOT") + drawdata.DefaultAttributeFontFile,
			expectedHeight: 15,
			expectedWidth:  58,
			hasError:       false,
		},
		{
			name:           "size invalid",
			str:            "Hello, World!",
			size:           0,
			fontFile:       os.Getenv("APP_ROOT") + drawdata.DefaultAttributeFontFile,
			expectedHeight: 0,
			expectedWidth:  0,
			hasError:       true,
		},
		{
			name:           "use default font",
			str:            "Hello, World!",
			size:           12,
			fontFile:       "",
			expectedHeight: 15,
			expectedWidth:  58,
			hasError:       false,
		},
		{
			name:     "bad size",
			str:      "Hello, World!",
			size:     0,
			fontFile: "",
			hasError: true,
		},
		{
			name:     "bad fontFile",
			str:      "Hello, World!",
			fontFile: "notExist.ttf",
			size:     12,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			height, width, err := GetTextSize(tt.str, tt.size, tt.fontFile)
			if (err != nil) != tt.hasError {
				t.Errorf("unexpected error: %v", err)
			}
			if height != tt.expectedHeight {
				t.Errorf("expected height %v, got %v", tt.expectedHeight, height)
			}
			if width != tt.expectedWidth {
				t.Errorf("expected width %v, got %v", tt.expectedWidth, width)
			}
		})
	}
}

func TestLoadFont(t *testing.T) {
	tests := []struct {
		name     string
		fontFile string
		hasError bool
	}{
		{
			name:     "correct",
			fontFile: os.Getenv("APP_ROOT") + drawdata.DefaultAttributeFontFile,
			hasError: false,
		},
		{
			name:     "file not found",
			fontFile: "./notExist.ttf",
			hasError: true,
		},
		{
			name:     "Bad font file",
			fontFile: "./dummy.ttf",
			hasError: true,
		},
	}
	dummy, err := os.Create("./dummy.ttf")
	assert.NoError(t, err)
	defer dummy.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadFont(tt.fontFile)
			if (err != nil) != tt.hasError {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
