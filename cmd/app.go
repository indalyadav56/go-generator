package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/indalyadav56/go-generator/file"
	"github.com/spf13/cobra"
)

// appCmd represents the app command
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
			fmt.Println("Please provide the name of the app")
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
	// tmpl, err := template.ParseGlob("templates/**/*.tmpl")
	// if err != nil {
	// 	panic(err)
	// }

	patterns := []string{
		"templates/*.tmpl",
		"templates/**/*.tmpl",
	}

	var allFiles []string
	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			panic(err)
		}
		allFiles = append(allFiles, files...)
	}

	// Parse the templates
	tmpl, err := template.ParseFiles(allFiles...)
	if err != nil {
		panic(err)
	}

	dirData := AddApp(appName)
	err = file.CreateStructure(dirPath, dirData, tmpl, appName)
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}

	if err := runSwaggerInit(dirPath); err != nil {
		log.Fatalf("Failed to run swag init: %v", err)
	}

}

func AddApp(title string) file.DirectoryStructure {
	structure := file.DirectoryStructure{
		// fmt.Sprintf("%s/constants", title):   {"constants.go"},
		// fmt.Sprintf("%s/routes", title):      {"routes.go"},
		// fmt.Sprintf("%s/dto", title):         {"dto.go"},
		// fmt.Sprintf("%s/models", title):      {"model.go"},
		fmt.Sprintf("%s/services", title): {"service.go"},
		// fmt.Sprintf("%s/services", title):    {"service.go", "service_test.go"},
		// fmt.Sprintf("%s/repository", title):  {"repository.go", "repository_test.go"},
		// fmt.Sprintf("%s/controllers", title): {"controller.go", "controller_test.go"},
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

	fmt.Printf("Executing command: swag %s\n", strings.Join(args, " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error initializing swagger: %v", err)
		log.Printf("Command output: %s", string(output))
		return err
	}

	fmt.Printf("Swagger initialization successful. Output: %s\n", string(output))
	return nil
}
