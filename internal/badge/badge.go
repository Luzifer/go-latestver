// Package badge contains a SVG-badge generator creating a badge from
// title, text and color
package badge

import (
	"bytes"
	"html/template"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"
)

const (
	xSpacing = 8
)

const (
	colorNameBlue        = "blue"
	colorNameBrightGreen = "brightgreen"
	colorNameGray        = "gray"
	colorNameGreen       = "green"
	colorNameLightGray   = "lightgray"
	colorNameOrange      = "orange"
	colorNameRed         = "red"
	colorNameYellow      = "yellow"
	colorNameYellowGreen = "yellowgreen"
)

var colorList = map[string]string{
	colorNameBlue:        "007ec6",
	colorNameBrightGreen: "4c1",
	colorNameGray:        "555",
	colorNameGreen:       "97CA00",
	colorNameLightGray:   "9f9f9f",
	colorNameOrange:      "fe7d37",
	colorNameRed:         "e05d44",
	colorNameYellow:      "dfb317",
	colorNameYellowGreen: "a4a61d",
}

// Create renders the badge and returns the SVG in minified but
// uncompressed form
func Create(title, text, color string) []byte {
	var buf bytes.Buffer

	titleW, _ := calculateTextWidth(title)
	textW, _ := calculateTextWidth(text)

	width := titleW + textW + 4*xSpacing //nolint:gomnd

	t, _ := assets.ReadFile("badge.svg.tpl")
	tpl, _ := template.New("svg").Parse(string(t))

	if c, ok := colorList[color]; ok {
		color = c
	}

	_ = tpl.Execute(&buf, map[string]any{
		"Width":       width,
		"TitleWidth":  titleW + 2*xSpacing,
		"Title":       title,
		"Text":        text,
		"TitleAnchor": titleW/2 + xSpacing,
		"TextAnchor":  titleW + textW/2 + 3*xSpacing,
		"Color":       color,
	})

	m := minify.New()
	m.AddFunc("image/svg+xml", svg.Minify)

	out, _ := m.Bytes("image/svg+xml", buf.Bytes())
	return out
}
