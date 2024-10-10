package file

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/indalyadav56/go-generator/format"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type DirectoryStructure map[string][]string

func CreateStructure(basePath string, structure DirectoryStructure, temp *template.Template) error {
	for dir, files := range structure {
		fullDirPath := filepath.Join(basePath, dir)

		err := CreateFolder(fullDirPath)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %w", fullDirPath, err)
		}

		// Check if there are any files to create, including empty strings
		if len(files) == 0 || (len(files) == 1 && files[0] == "") {
			fmt.Printf("Empty directory created: %s\n", fullDirPath)
			continue
		}

		for _, file := range files {
			if file == "" {
				continue
			}

			contentData, _ := ParseContent(temp, file, dir, basePath)
			err = CreateFile(filepath.Join(fullDirPath, file), contentData)
			if err != nil {
				log.Fatalf("Failed to create file: %v", err)
			}
		}
	}
	return nil
}

func CreateFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}
	return nil
}

func CreateFile(filePath, content string) error {
	folder := filepath.Dir(filePath)
	if err := CreateFolder(folder); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	if content != "" {
		_, err = file.WriteString(content)
		if err != nil {
			return fmt.Errorf("error writing content to file: %v", err)
		}
	}

	return nil
}

// ParseContent processes a Go template and returns the formatted Go code.
func ParseContent(tmpl *template.Template, fileName, dir, projectTitle string) (string, error) {
	caser := cases.Title(language.English)
	capitalized := caser.String(projectTitle)

	data := map[string]string{
		"IServiceName": capitalized,
		"ServiceName":  projectTitle,
		"AppName":      projectTitle,
	}

	if strings.Contains(fileName, "app.log") {
		return "", nil
	}

	if strings.Contains(fileName, "auth") {
		data["IServiceName"] = "Auth"
		data["ServiceName"] = projectTitle
		data["AppName"] = "auth"
	}

	if strings.Contains(fileName, "user") {
		data["IServiceName"] = "User"
		data["ServiceName"] = projectTitle
		data["AppName"] = "user"
	}

	templateName, isFormat := getTemplateName(fileName, dir)

	var output bytes.Buffer

	err := tmpl.ExecuteTemplate(&output, templateName, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	if isFormat {
		formattedOutput, err := format.FormatGoCode(output.Bytes())
		if err != nil {
			return "", fmt.Errorf("failed to format Go code: %w", err)
		}
		return string(formattedOutput), nil
	}

	return output.String(), nil
}

func getTemplateName(fileName, dir string) (templateName string, isFormat bool) {
	isFormat = true

	patterns := map[string]string{
		"main":              "main",
		"service_test":      "service_test",
		"dto_test":          "dto_test",
		"controller_test":   "controller_test",
		"auth_service":      "auth_service",
		"service":           "service",
		"repository_test":   "repository_test",
		"repository":        "repository",
		"constant":          "constant",
		"auth_controller":   "auth_controller",
		"controller":        "controller",
		"logger_middleware": "logger_middleware",
		"auth_middleware":   "auth_middleware",
		"routes":            "routes",
		"dto":               "dto",
		"model":             "model",
	}

	nonFormattedPatterns := map[string]string{
		"makefile":   "makefile",
		"dockerfile": "dockerfile",
		"readme":     "readme",
	}

	// Check if the fileName matches any of the non-format cases
	for key, template := range nonFormattedPatterns {
		if strings.Contains(strings.ToLower(fileName), key) {
			return template, false
		}
	}

	for pattern, template := range patterns {
		if strings.Contains(fileName, pattern) {
			return template, true
		}
	}

	if strings.Contains(dir, "database") {
		return "database", true
	} else if strings.Contains(dir, "config") {
		return "config", true
	}

	return "unknown", false
}
