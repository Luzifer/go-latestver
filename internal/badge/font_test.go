package badge

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedFontHash(t *testing.T) {
	// Check the embedded font did not change
	font, err := assets.ReadFile("DejaVuSans.ttf")
	require.NoError(t, err)

	assert.Equal(t,
		"3fdf69cabf06049ea70a00b5919340e2ce1e6d02b0cc3c4b44fb6801bd1e0d22",
		fmt.Sprintf("%x", sha256.Sum256(font)))
}

func TestStringLength(t *testing.T) {
	// As the font is embedded into the source the length calculation should not change
	w, err := calculateTextWidth("Test 123 öäüß … !@#%&")
	require.NoError(t, err)
	assert.Equal(t, 138, w)
}
