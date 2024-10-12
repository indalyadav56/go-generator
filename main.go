package main

import (
	"fmt"
	"log"
	"text/template"

	"github.com/indalyadav56/go-generator/cmd"
	"github.com/indalyadav56/go-generator/file"
)

func main() {
	cmd.Execute()
	// // CreateProject("/todo/internal")

	// tmpl, err := template.ParseGlob("templates/*.tmpl")
	// if err != nil {
	// 	panic(err)
	// }

	// dirData := AddApp("task")
	// err = file.CreateStructure("./todo/internal", dirData, tmpl)
	// if err != nil {
	// 	log.Fatalf("Failed to create structure: %v\n", err)
	// }

}

func CreateProject(projectTitle string) {
	tmpl, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	structure := file.DirectoryStructure{
		"auth/routes": {"auth_routes.go"},
	}

	err = file.CreateStructure(projectTitle, structure, tmpl)
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}

}

func AddApp(title string) file.DirectoryStructure {
	return file.DirectoryStructure{
		fmt.Sprintf("%s/constants", title):   {"constant.go"},
		fmt.Sprintf("%s/routes", title):      {"routes.go"},
		fmt.Sprintf("%s/dto", title):         {fmt.Sprintf("%s_dto.go", title)},
		fmt.Sprintf("%s/models", title):      {fmt.Sprintf("%s_model.go", title)},
		fmt.Sprintf("%s/services", title):    {fmt.Sprintf("%s_service.go", title)},
		fmt.Sprintf("%s/repository", title):  {fmt.Sprintf("%s_repository.go", title)},
		fmt.Sprintf("%s/controllers", title): {fmt.Sprintf("%s_controller.go", title)},
	}
}
