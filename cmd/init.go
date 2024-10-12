package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/indalyadav56/go-generator/file"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "This command initializes a golang project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args", args)
		apps, _ := cmd.Flags().GetStringSlice("app")
		fmt.Println("apps", apps)

		if len(args) == 0 {
			fmt.Println("provide the name of the project")
			return
		}
		name := args[0]
		CreateProject(strings.ToLower(name))

		// create apps
		for _, v := range apps {
			CreateApp(strings.ToLower(v), fmt.Sprintf("./%s/internal", name))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringSliceP("app", "", []string{}, "app names to initialize")
}

func CreateProject(projectTitle string) {
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
		"internal":                          {""},
		".":                                 {".gitignore", "README.md", "Dockerfile", "docker-compose.yml", "Makefile", ".env"},
	}

	err = file.CreateStructure(projectTitle, structure, tmpl, "")
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}

	initGoModule(projectTitle)
	copyCommonPkg(projectTitle)

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
