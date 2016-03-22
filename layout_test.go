package formats

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLayout(t *testing.T) {

	// XXX
	file, err := os.Open("samples/arj/tiny.arj")
	defer file.Close()
	assert.Equal(t, nil, err)

	layout, err := ParseLayout(file)
	assert.Equal(t, nil, err)
	fmt.Println(layout) // XXX test layout
}
