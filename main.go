package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/indalyadav56/go-generator/file"
)

func main() {
	var projectTitle string

	fmt.Print("Enter your project title: ")
	fmt.Scanln(&projectTitle)

	tmpl, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	structure := file.DirectoryStructure{
		fmt.Sprintf("cmd/%s", projectTitle): {"main.go"},
		"config":                            {"env.go"},
		"database":                          {"postgres.go"},
		"docs":                              {""},

		// auth
		"internal/auth/constants":   {"constant.go"},
		"internal/auth/routes":      {"routes.go"},
		"internal/auth/dto":         {"auth_dto.go"},
		"internal/auth/services":    {"auth_service.go"},
		"internal/auth/repository":  {"auth_repository.go"},
		"internal/auth/controllers": {"auth_controller.go"},

		// user
		"internal/user/constants":   {"constant.go"},
		"internal/user/routes":      {"routes.go"},
		"internal/user/dto":         {"user_dto.go"},
		"internal/user/models":      {"user_model.go"},
		"internal/user/services":    {"user_service.go"},
		"internal/user/repository":  {"user_repository.go"},
		"internal/user/controllers": {"user_controller.go"},

		// todo
		"internal/todo/constants":   {"constant.go"},
		"internal/todo/routes":      {"routes.go"},
		"internal/todo/dto":         {fmt.Sprintf("%s_dto.go", projectTitle)},
		"internal/todo/models":      {fmt.Sprintf("%s_model.go", projectTitle)},
		"internal/todo/services":    {fmt.Sprintf("%s_service.go", projectTitle)},
		"internal/todo/repository":  {fmt.Sprintf("%s_repository.go", projectTitle)},
		"internal/todo/controllers": {fmt.Sprintf("%s_controller.go", projectTitle)},

		//
		"logs": {"app.log"},
		".":    {".gitignore", "README.md", "Dockerfile", "docker-compose.yml", "Makefile", ".env"},

		// "scripts":                           {"setup.sh", "deploy.sh"},
		// "web/static":       {""},
		// "web/templates":    {""},
		// "test/unit":        {""},
		// "test/integration": {""},
	}

	// Call the createStructure function to create the directories and files
	err = file.CreateStructure(projectTitle, structure, tmpl)
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}

	initGoModule(projectTitle)
	copyCommonPkg(projectTitle)

	// Run 'go mod tidy' to install dependencies and create the go.sum file
	err = runGoModTidy("./" + projectTitle)
	if err != nil {
		log.Fatalf("Failed to run 'go mod tidy': %v", err)
	}

}

// runGoModTidy runs 'go mod tidy' to resolve dependencies and generate the go.sum file
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

func initGoModule(projectTitle string) error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Create the "custom" folder within the current directory
	customDir := filepath.Join(currentDir, projectTitle)

	// Ensure the "custom" folder exists
	err = os.MkdirAll(customDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating 'custom' folder: %v", err)
	}

	// Run 'go mod init' in the "custom" folder
	cmd := exec.Command("go", "mod", "init", projectTitle)
	cmd.Dir = customDir
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error initializing Go module: %v", err)
	}

	fmt.Println("Initialized Go module in:", customDir)
	fmt.Println("Go module initialized:", projectTitle)
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

// CopyFolder recursively copies a folder and its files from src to dst.
func CopyFolder(src, dst string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create the target path by replacing the source base path with the destination base path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		// If it's a directory, create the corresponding directory at the destination
		if info.IsDir() {
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			return nil
		}

		// If it's a file, copy it to the target directory
		err = copyFile(path, targetPath)
		if err != nil {
			return fmt.Errorf("failed to copy file: %v", err)
		}

		return nil
	})
	return err
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the file contents from source to destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Ensure the destination file has the same permissions as the source
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, sourceInfo.Mode())
	if err != nil {
		return err
	}

	fmt.Printf("Copied file: %s to %s\n", src, dst)
	return nil
}
