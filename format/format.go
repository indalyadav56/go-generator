package format

import (
	"go/format"
)

// format Go source code
func FormatGoCode(input []byte) ([]byte, error) {
	formattedOutput, err := format.Source(input)
	if err != nil {
		return nil, err
	}
	return formattedOutput, nil
}
