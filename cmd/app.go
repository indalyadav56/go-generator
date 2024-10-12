package cmd

import (
	"fmt"
	"log"
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
		CreateApp(appName, dirPath)

	},
}

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.Flags().StringP("name", "n", "", "name of the app")
}

func CreateApp(appName, dirPath string) {
	tmpl, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	dirData := AddApp(appName)
	err = file.CreateStructure(dirPath, dirData, tmpl, appName)
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}

}

func AddApp(title string) file.DirectoryStructure {
	structure := file.DirectoryStructure{
		fmt.Sprintf("%s/constants", title):   {"constant.go"},
		fmt.Sprintf("%s/routes", title):      {"routes.go"},
		fmt.Sprintf("%s/dto", title):         {fmt.Sprintf("%s_dto.go", title)},
		fmt.Sprintf("%s/models", title):      {fmt.Sprintf("%s_model.go", title)},
		fmt.Sprintf("%s/services", title):    {fmt.Sprintf("%s_service.go", title)},
		fmt.Sprintf("%s/repository", title):  {fmt.Sprintf("%s_repository.go", title)},
		fmt.Sprintf("%s/controllers", title): {fmt.Sprintf("%s_controller.go", title)},
	}

	if title == "auth" || title == "authentication" {
		delete(structure, fmt.Sprintf("%s/repository", title))
	}
	return structure
}
