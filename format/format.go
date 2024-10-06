package format

import (
	"fmt"
	"go/format"
)

// Function to format Go source code and handle errors
func FormatGoCode(input []byte) ([]byte, error) {
	// Use go/format to format the source code
	formattedOutput, err := format.Source(input)
	if err != nil {
		// Handle formatting errors (usually due to invalid Go syntax in the template)
		fmt.Println("Error formatting Go code:", err)
		fmt.Println("Unformatted output:", string(input)) // Print unformatted output for debugging
		return nil, err
	}
	return formattedOutput, nil
}
