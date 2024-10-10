package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/indalyadav56/go-generator/cmd"
	"github.com/indalyadav56/go-generator/file"
)

func main() {
	cmd.Execute()
	return

	var projectTitle string

	projectTitle = "todo"

	// fmt.Print("Enter your project title: ")
	// fmt.Scanln(&projectTitle)

	tmpl, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	structure := file.DirectoryStructure{
		fmt.Sprintf("cmd/%s", projectTitle): {"main.go"},
		"config":                            {"env.go"},
		"database":                          {"postgres.go"},
		"middlewares":                       {"logger_middleware.go", "auth_middleware.go"},
		"docs":                              {""},

		// // auth
		// "internal/auth/constants":   {"constant.go"},
		// "internal/auth/routes":      {"auth_routes.go"},
		// "internal/auth/dto":         {"auth_dto.go"},
		// "internal/auth/services":    {"auth_service.go"},
		// "internal/auth/controllers": {"auth_controller.go"},

		// // user
		// "internal/user/constants":   {"constant.go"},
		// "internal/user/routes":      {"user_routes.go"},
		// "internal/user/dto":         {"user_dto.go"},
		// "internal/user/models":      {"user_model.go"},
		// "internal/user/services":    {"user_service.go"},
		// "internal/user/repository":  {"user_repository.go"},
		// "internal/user/controllers": {"user_controller.go"},

		// "internal/todo/constants":   {"constant.go"},
		// "internal/todo/routes":      {fmt.Sprintf("%s_routes.go", projectTitle)},
		// "internal/todo/dto":         {fmt.Sprintf("%s_dto.go", projectTitle)},
		// "internal/todo/models":      {fmt.Sprintf("%s_model.go", projectTitle)},
		// "internal/todo/services":    {fmt.Sprintf("%s_service.go", projectTitle)},
		// "internal/todo/repository":  {fmt.Sprintf("%s_repository.go", projectTitle)},
		// "internal/todo/controllers": {fmt.Sprintf("%s_controller.go", projectTitle)},

		//
		"logs": {"app.log"},
		".":    {".gitignore", "README.md", "Dockerfile", "docker-compose.yml", "Makefile", ".env"},

		// "scripts":                           {"setup.sh", "deploy.sh"},
		// "web/static":       {""},
		// "web/templates":    {""},
		// "test/unit":        {""},
		// "test/integration": {""},
	}

	err = file.CreateStructure(projectTitle, structure, tmpl)
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}

	// err = file.CreateStructure(projectTitle, AddApp("task"), tmpl)
	// if err != nil {
	// 	log.Fatalf("Failed to create structure: %v\n", err)
	// }

	initGoModule(projectTitle)
	copyCommonPkg(projectTitle)
	runSwagInit(fmt.Sprintf("%s/cmd/%s/main.go", projectTitle, projectTitle))

	err = runGoModTidy("./" + projectTitle)
	if err != nil {
		log.Fatalf("Failed to run 'go mod tidy': %v", err)
	}

}

func runGoModTidy(basePath string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = basePath
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running 'go mod tidy': %v", err)
	}
	fmt.Println("'go mod tidy' executed successfully, dependencies resolved.")
	return nil
}

func runSwagInit(basePath string) error {
	cmd := exec.Command("swag", "init")
	cmd.Dir = basePath
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running 'go mod tidy': %v", err)
	}
	fmt.Println("'go mod tidy' executed successfully, dependencies resolved.")
	return nil
}

func initGoModule(projectTitle string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	customDir := filepath.Join(currentDir, projectTitle)

	err = os.MkdirAll(customDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating 'custom' folder: %v", err)
	}

	cmd := exec.Command("go", "mod", "init", projectTitle)
	cmd.Dir = customDir
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error initializing Go module: %v", err)
	}

	return nil
}

func copyCommonPkg(projectTitle string) {
	srcDir := "./pkg"

	err := CopyFolder(srcDir, projectTitle+"/pkg")
	if err != nil {
		log.Fatalf("Failed to copy folder: %v", err)
	}

	fmt.Println("Folder copied successfully!")
}

func CopyFolder(src, dst string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			return nil
		}

		err = copyFile(path, targetPath)
		if err != nil {
			return fmt.Errorf("failed to copy file: %v", err)
		}

		return nil
	})
	return err
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, sourceInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

func AddApp(title string) file.DirectoryStructure {
	return file.DirectoryStructure{
		fmt.Sprintf("internal/%s/constants", title):   {"constant.go"},
		fmt.Sprintf("internal/%s/routes", title):      {"routes.go"},
		fmt.Sprintf("internal/%s/dto", title):         {fmt.Sprintf("%s_dto.go", title)},
		fmt.Sprintf("internal/%s/models", title):      {fmt.Sprintf("%s_model.go", title)},
		fmt.Sprintf("internal/%s/services", title):    {fmt.Sprintf("%s_service.go", title)},
		fmt.Sprintf("internal/%s/repository", title):  {fmt.Sprintf("%s_repository.go", title)},
		fmt.Sprintf("internal/%s/controllers", title): {fmt.Sprintf("%s_controller.go", title)},
	}
}
