package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/indalyadav56/go-generator/file"
	"github.com/indalyadav56/go-generator/templates"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "to create a new app",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var dirPath string

		if len(args) > 0 {
			dirPath = args[0]
		}

		appName, _ := cmd.Flags().GetString("name")
		if appName == "" {
			return
		}
		// apiFramework, _ := cmd.Flags().GetString("framework")
		CreateApp(appName, dirPath)

	},
}

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.Flags().StringP("name", "n", "", "name of the app")
}

func CreateApp(appName, dirPath string) {
	// Parse all templates from the embedded filesystem
	tmpl := template.New("")
	
	err := fs.WalkDir(templates.TemplateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			content, err := templates.TemplateFS.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = tmpl.New(filepath.Base(path)).Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})
	
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	dirData := AddApp(appName)
	err = file.CreateStructure(dirPath, dirData, tmpl, appName)
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}
}

func AddApp(title string) file.DirectoryStructure {
	structure := file.DirectoryStructure{
		fmt.Sprintf("%s/constants", title):   {"constants.go"},
		fmt.Sprintf("%s/routes", title):      {"routes.go"},
		fmt.Sprintf("%s/dto", title):         {"dto.go"},
		fmt.Sprintf("%s/models", title):      {"model.go"},
		fmt.Sprintf("%s/services", title):    {"service.go"},
		fmt.Sprintf("%s/services", title):    {"service.go", "service_test.go"},
		fmt.Sprintf("%s/repository", title):  {"repository.go", "repository_test.go"},
		fmt.Sprintf("%s/controllers", title): {"controller.go", "controller_test.go"},
	}

	if title == "auth" || title == "authentication" {
		delete(structure, fmt.Sprintf("%s/repository", title))
		delete(structure, fmt.Sprintf("%s/models", title))
	}
	return structure
}

func runSwaggerInit(dirPath string) error {
	var projectTitle string

	if strings.Contains(dirPath, "/") {
		data := strings.Split(dirPath, "/")
		projectTitle = data[1]
	}

	// dirPath = "./backend/cmd/backend"
	dirPath = fmt.Sprintf("./%s/cmd/%s", projectTitle, projectTitle)
	args := []string{"init", "-o", "../../docs", fmt.Sprintf("./cmd/%s", projectTitle), dirPath}
	cmd := exec.Command("swag", args...)
	cmd.Dir = dirPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error initializing swagger: %v", err)
		log.Printf("Command output: %s", string(output))
		return err
	}

	return nil
}
