package cmd

import (
	"fmt"
	"io/fs"
	"log"
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
		CreateApp(appName, dirPath)

	},
}

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.Flags().StringP("name", "n", "", "name of the app")
}

func CreateApp(appName, dirPath string) {
	defer runGoModTidy(dirPath)

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

	dirData := AddApp(appName, dirPath)
	templateData := map[string]interface{}{
		"AppName": appName,
	}
	err = file.CreateStructure(dirPath, dirData, tmpl, appName, templateData)
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}
}

func AddApp(title, dirPath string) file.DirectoryStructure {
	structure := file.DirectoryStructure{
		fmt.Sprintf("%s/constants", title):  {"constants.go"},
		fmt.Sprintf("%s/routes", title):     {"routes.go"},
		fmt.Sprintf("%s/dto", title):        {"dto.go"},
		fmt.Sprintf("%s/models", title):     {"model.go"},
		fmt.Sprintf("%s/services", title):   {"service.go"},
		fmt.Sprintf("%s/repository", title): {"repository.go"},
		fmt.Sprintf("%s/handlers", title):   {"handler.go"},
	}

	if title == "auth" || title == "authentication" {
		delete(structure, fmt.Sprintf("%s/repository", title))
		delete(structure, fmt.Sprintf("%s/models", title))
	}
	return structure
}
