package file

import (
	"testing"
)

func TestGetTemplateName(t *testing.T) {
	testCases := []struct {
		fileName, dir string
	}{
		{"auth/services/auth_service.go", ""},
		{"user/controllers/user_controller_test.go", ""},
		{"Makefile", ""},
		{"config/config.go", "config"},
		{"database/connection.go", "database"},
		{"unknown_file.go", ""},
	}

	for _, tc := range testCases {
		template, isFormat := getTemplateName(tc.fileName, tc.dir)
		println(tc.fileName, "->", template, isFormat)
	}
}
