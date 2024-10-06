package format

import (
	"fmt"
	"go/format"
)

// format Go source code
func FormatGoCode(input []byte) ([]byte, error) {
	formattedOutput, err := format.Source(input)
	if err != nil {
		fmt.Println("Error formatting Go code:", err)
		fmt.Println("Unformatted output:", string(input))
		return nil, err
	}
	return formattedOutput, nil
}
