package utils

import (
	_ "embed"
)

type FontConfig struct {
	FontFile  []byte
	FontScale float64
}

//go:embed fonts/LXGWWenKaiTC-EmojiCompletion.ttf
var fontFile []byte

var BuiltinFont = FontConfig{
	FontFile:  fontFile,
	FontScale: 0.9,
}
