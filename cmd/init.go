package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/indalyadav56/go-generator/file"
	"github.com/indalyadav56/go-generator/templates"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "This command initializes a golang project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		apps, _ := cmd.Flags().GetStringSlice("app")
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
	// Debug: List all files in embedded FS
	entries, err := templates.TemplateFS.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subEntries, err := templates.TemplateFS.ReadDir(entry.Name())
			if err != nil {
				panic(err)
			}
			for _, subEntry := range subEntries {
				fmt.Printf("  - %s/%s\n", entry.Name(), subEntry.Name())
			}
		}
	}

	// Parse the templates using the embedded file system
	tmpl, err := template.ParseFS(templates.TemplateFS,
		"templates/app/app_config.tmpl",
		"templates/app/app_config_router.tmpl",
		"templates/app/gorm_app_config.tmpl",
		"templates/config/config.tmpl",
		"templates/docker/dockerfile.tmpl",
		"templates/docker/docker-compose.tmpl",
		"templates/env/env.tmpl",
		"templates/makefile/makefile.tmpl",
		"templates/readme/readme.tmpl",
		"templates/models/model.tmpl",
		"templates/models/model_test.tmpl",
		"templates/main.tmpl",
		"templates/gin/routes.tmpl",
		"templates/gin/controller.tmpl",
		"templates/gin/auth_middleware.tmpl",
		"templates/gin/logger_middleware.tmpl",
		"templates/gin/app_config_router.tmpl",
		"templates/constants/constant.tmpl",
		"templates/gorm/postgres_db.tmpl",
		"templates/gorm/db_logger.tmpl",
	)
	if err != nil {
		panic(err)
	}

	structure := file.DirectoryStructure{
		fmt.Sprintf("cmd/%s", projectTitle): {"main.go"},
		"config":                            {"config.go"},
		// "db":                                {"postgres.go"},
		// "migrations":                        {""},
		// "docs":                              {""},
		"internal/app": {"app.go", "deps.go"},
		".":            {".gitignore", "README.md", "Dockerfile", "docker-compose.yml", "Makefile", ".env"},
	}

	err = file.CreateStructure(projectTitle, structure, tmpl, "")
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}

	initGoModule(projectTitle)

	// wg := new(sync.WaitGroup)

	// wg.Add(1)
	// go downloadPkgFromGithub(projectTitle, wg)

	// err = runGoModTidy("./" + projectTitle)
	// if err != nil {
	// 	log.Fatalf("Failed to run 'go mod tidy': %v", err)
	// }

	// wg.Wait()
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
	fmt.Println("projectTitle", projectTitle)

	fmt.Println("currentDir", currentDir)

	customDir := filepath.Join(currentDir, projectTitle)
	fmt.Println("customDir", customDir)

	err = os.MkdirAll(customDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	cmd := exec.Command("go", "mod", "init", projectTitle)
	cmd.Dir = customDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		// Try to clean up if initialization fails
		os.RemoveAll(customDir)
		log.Fatalf("Error initializing Go module: %v", err)
	}

	// Create go.mod file with required dependencies
	goModContent := `module ` + projectTitle + `

go 1.21

require (
	github.com/go-chi/chi/v5 v5.0.11
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/jackc/pgx/v5 v5.5.1
	golang.org/x/crypto v0.17.0
)
`
	err = os.WriteFile(filepath.Join(customDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		log.Fatalf("Error writing go.mod file: %v", err)
	}

	return nil
}

func downloadPkgFromGithub(projectTitle string, wg *sync.WaitGroup) {
	defer wg.Done()
	repoURL := "https://github.com/indalyadav56/go-generator"
	pkgPath := "pkg"
	targetDir := filepath.Join(projectTitle, "pkg")

	// Remove existing pkg directory if it exists
	if err := os.RemoveAll(targetDir); err != nil {
		log.Fatalf("Error removing existing pkg directory: %v", err)
	}

	// Create the target directory
	err := os.MkdirAll(filepath.Dir(targetDir), 0755)
	if err != nil {
		log.Fatalf("Error creating target directory: %v", err)
	}

	// Run git sparse-checkout to download only the pkg directory
	cmd := exec.Command("git", "clone", "--depth", "1", "--filter=blob:none", "--sparse", repoURL, targetDir+"_temp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error cloning repository: %v", err)
	}

	// Change to the temp directory
	tempDir := targetDir + "_temp"
	cmd = exec.Command("git", "sparse-checkout", "set", pkgPath)
	cmd.Dir = tempDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error setting sparse-checkout: %v", err)
	}

	// Move the pkg directory to the target location
	srcPath := filepath.Join(tempDir, pkgPath)
	err = os.Rename(srcPath, targetDir)
	if err != nil {
		log.Fatalf("Error moving pkg directory: %v", err)
	}

	// Clean up the temp directory
	err = os.RemoveAll(tempDir)
	if err != nil {
		log.Printf("Warning: Error cleaning up temp directory: %v", err)
	}

	fmt.Println("Successfully downloaded pkg directory from GitHub!")
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
