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
	tmpl, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	projectTitle := "todo"

	structure := file.DirectoryStructure{
		fmt.Sprintf("cmd/%s", projectTitle): {fmt.Sprintf("%s.go", projectTitle)},
		"config":                            {"env.go"},
		"database":                          {"postgres.go"},
		"repository":                        {fmt.Sprintf("%s_repository.go", projectTitle)},
		"controllers":                       {fmt.Sprintf("%s_controller.go", projectTitle)},
		"services":                          {fmt.Sprintf("%s_service.go", projectTitle)},
		"logs":                              {"app.log"},
		".":                                 {"README.md"},
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

// Function to initialize a Go module (go mod init)
// initGoModule initializes a Go module in the "custom" folder within the current working directory.
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
