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
type templatePattern struct {
	pattern  string
	template string
	isFormat bool
}

var templatePatterns = []templatePattern{
	{"makefile", "makefile", false},
	{"dockerfile", "dockerfile", false},
	{"readme", "readme", false},

	// Specific patterns
	{"auth_constant", "auth_constant", true},
	{"auth_service_test", "auth_service_test", true},
	{"auth_controller_test", "controller_test", true},
	{"auth_integration_test", "controller_test", true},

	{"auth_dto", "auth_dto", true},
	{"auth_routes", "auth_routes", true},
	{"auth_service", "auth_service", true},
	{"auth_controller", "auth_controller", true},
	{"logger_middleware", "logger_middleware", true},
	{"auth_middleware", "auth_middleware", true},

	{"gitignore", "gitignore", false},
	{"docker-compose", "compose", false},

	// General patterns
	{"env.go", "config", true},
	{"app.go", "app_config", true},
	{"router.go", "app_config_router", true},

	{"controller_test", "controller_test", true},
	{"integration_test", "controller_test", true},
	{"service_test", "service_test", true},
	{"repository_test", "repository_test", true},
	// {"model_test", "model_test", true},

	{"service", "service", true},
	{"controller", "controller", true},
	{"repository", "repository", true},
	{"routes", "routes", true},
	{"dto", "dto", true},
	{"model", "model", true},
	{"constant", "constant", true},

	{"db_logger", "db_logger", true},
	{"postgres", "postgres_db", true},
	{"main", "main", true},
	{"env", "env", false},
}

func CreateStructure(basePath string, structure DirectoryStructure, temp *template.Template, appName string) error {
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

			contentData, _ := ParseContent(temp, file, dir, basePath, appName)
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
func ParseContent(tmpl *template.Template, fileName, dir, projectTitle, appName string) (string, error) {

	if strings.Contains(projectTitle, ".") {
		data := strings.Split(projectTitle, "/")
		projectTitle = data[1]
	}

	caser := cases.Title(language.English)
	capitalized := caser.String(projectTitle)
	appCapitalized := caser.String(appName)

	data := map[string]string{
		"IServiceName": capitalized,
		"ServiceName":  projectTitle,
		"IAppName":     appCapitalized,
		"AppName":      appName,
	}

	if strings.Contains(fileName, "app.log") {
		return "", nil
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
	lowerFileName := strings.ToLower(fileName)
	baseName := filepath.Base(lowerFileName)

	for _, tp := range templatePatterns {
		if strings.Contains(baseName, tp.pattern) {
			return tp.template, tp.isFormat
		}
	}

	if strings.HasSuffix(baseName, "_test.go") {
		return filepath.Base(strings.TrimSuffix(baseName, "_test.go")) + "_test", true
	}

	return "unknown", false
}
