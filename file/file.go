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
		// Append the base path to the directory
		fullDirPath := filepath.Join(basePath, dir)

		// Create the directory and any necessary parents
		err := CreateFolder(fullDirPath)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %w", fullDirPath, err)
		}

		// Check if there are any files to create, including empty strings
		if len(files) == 0 || (len(files) == 1 && files[0] == "") {
			// No files to create, just create the directory
			fmt.Printf("Empty directory created: %s\n", fullDirPath)
			continue
		}

		// Create each file in the directory
		for _, file := range files {
			if file == "" {
				continue
			}

			// Read tmpl file and parse data
			contentData, _ := ParseContent(temp, file, dir, basePath)
			// Create file with content
			err = CreateFile(filepath.Join(fullDirPath, file), contentData)
			if err != nil {
				log.Fatalf("Failed to create file: %v", err)
			}
		}
	}
	return nil
}

// CreateFolder creates the specified folder if it doesn't exist.
func CreateFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}
	fmt.Println("Created folder:", folderPath)
	return nil
}

// CreateFile creates a file at the specified filePath and writes content to it.
func CreateFile(filePath, content string) error {
	// Ensure the folder for the file exists
	folder := filepath.Dir(filePath)
	if err := CreateFolder(folder); err != nil {
		return err
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write content if any
	if content != "" {
		_, err = file.WriteString(content)
		if err != nil {
			return fmt.Errorf("error writing content to file: %v", err)
		}
	}

	fmt.Println("Created file:", filePath)
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

	templateName := "main"
	isFormat := true

	switch {
	case strings.Contains(fileName, "main"):
		templateName = "main"

	case strings.Contains(fileName, "service_test"):
		templateName = "service_test"

	case strings.Contains(fileName, "dto_test"):
		templateName = "dto_test"

	case strings.Contains(fileName, "controller_test"):
		templateName = "controller_test"

	case strings.Contains(fileName, "service"):
		templateName = "service"

	case strings.Contains(fileName, "repository_test"):
		templateName = "repository_test"

	case strings.Contains(fileName, "repository"):
		templateName = "repository"

	case strings.Contains(dir, "database"):
		templateName = "database"

	case strings.Contains(dir, "config"):
		templateName = "config"

	case strings.Contains(fileName, "constant"):
		templateName = "constant"

	case strings.Contains(fileName, "controller"):
		templateName = "controller"

	case strings.Contains(fileName, "logger_middleware"):
		templateName = "logger_middleware"

	case strings.Contains(fileName, "auth_middleware"):
		templateName = "auth_middleware"

	case strings.Contains(fileName, "routes"):
		templateName = "routes"

	case strings.Contains(fileName, "dto"):
		templateName = "dto"

	case strings.Contains(fileName, "model"):
		templateName = "model"

	case strings.Contains(strings.ToLower(fileName), "makefile"):
		templateName = "makefile"
		isFormat = false
	case strings.Contains(strings.ToLower(fileName), "dockerfile"):
		templateName = "dockerfile"
		isFormat = false
	case strings.Contains(strings.ToLower(fileName), "readme"):
		templateName = "readme"
		isFormat = false
	default:
		templateName = "unknown"
		isFormat = false
	}

	fmt.Println("templateName", templateName)

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
