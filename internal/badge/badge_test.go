package badge

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBadgeTemplate(t *testing.T) {
	raw, err := assets.ReadFile("badge.svg.tpl")
	require.NoError(t, err)

	_, err = template.New("svg").Parse(string(raw))
	require.NoError(t, err)
}

func TestRenderBadge(t *testing.T) {
	badge := Create("golang", "test", "green")
	assert.Equal(t,
		[]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="90" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><mask id="a"><rect width="90" height="20" rx="3" fill="#fff"/></mask><g mask="url(#a)"><path fill="#555" d="M0 0h53v20H0z"/><path fill="#97ca00" d="M53 0H90v20H53z"/><path fill="url(#b)" d="M0 0h90v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="26" y="15" fill="#010101" fill-opacity=".3">golang</text><text x="26" y="14">golang</text><text x="71" y="15" fill="#010101" fill-opacity=".3">test</text><text x="71" y="14">test</text></g></svg>`),
		badge)
}
