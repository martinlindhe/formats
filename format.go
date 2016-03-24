package formats

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// FormatDescription ...
type FormatDescription struct {
	Format Format `json:"format"`
}

// Format ...
type Format struct {
	Name    string   `json:"name"`
	Mime    string   `json:"mime"`
	Details []string `json:"details"`
}

// ReadFormatDescription ...
func ReadFormatDescription(formatName string) (*Format, error) {

	formatFile := "./formats/" + formatName + ".yml"

	if !exists(formatFile) {
		return nil, fmt.Errorf("Unknown format %s", formatFile)
	}

	data, err := ioutil.ReadFile(formatFile)
	if err != nil {
		return nil, err
	}

	desc := FormatDescription{}
	err = yaml.Unmarshal(data, &desc)
	return &desc.Format, err
}
