package file

import (
	"fmt"
	"testing"
)

func TestGetTemplateName(t *testing.T) {

	t.Run("get template name", func(t *testing.T) {
		template, isFormat := getTemplateName("routes/routes.go", "")
		fmt.Println("template:", template)
		fmt.Println("isFormat:", isFormat)
	})
}
