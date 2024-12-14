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
	"github.com/indalyadav56/go-generator/templates"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "This command initializes a golang project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		apps, _ := cmd.Flags().GetStringSlice("app")
		framework, _ := cmd.Flags().GetString("framework")
		frontend, _ := cmd.Flags().GetString("frontend")
		if len(args) == 0 {
			fmt.Println("provide the name of the project")
			return
		}
		name := args[0]
		CreateProject(strings.ToLower(name), framework, frontend)

		// create apps
		for _, v := range apps {
			CreateApp(strings.ToLower(v), fmt.Sprintf("./%s/internal", name))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringSliceP("app", "", []string{}, "app names to initialize")
	initCmd.Flags().StringP("framework", "", "", "web framework to use (e.g., htmx)")
	initCmd.Flags().StringP("frontend", "", "", "frontend framework to use (e.g., react, htmx)")
}

func CreateProject(projectTitle string, framework string, frontend string) {
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

	templatePaths := []string{
		"templates/app/app.tmpl",
		"templates/app/deps.tmpl",
		"templates/config/config.tmpl",
		"templates/docker/dockerfile.tmpl",
		"templates/docker/docker-compose.tmpl",
		"templates/env/env.tmpl",
		"templates/makefile/makefile.tmpl",
		"templates/readme/readme.tmpl",
		"templates/models/model.tmpl",
		"templates/main.tmpl",
		"templates/gin/routes.tmpl",
		"templates/gin/auth_middleware.tmpl",
		"templates/gin/logger_middleware.tmpl",
		"templates/constants/constant.tmpl",
	}

	if frontend == "htmx" {
		templatePaths = append(templatePaths,
			"templates/htmx/base.tmpl",
			"templates/htmx/index.tmpl",
			"templates/htmx/style.tmpl",
		)
	}

	tmpl, err := template.ParseFS(templates.TemplateFS, templatePaths...)
	if err != nil {
		panic(err)
	}

	structure := file.DirectoryStructure{
		fmt.Sprintf("cmd/%s", "api"): {"main.go"},
		"config":                     {"config.go"},
		"docs":                       {""},
		"scripts":                    {"build.sh"},
		"logs":                       {""},
		"internal/app":               {"app.go", "deps.go"},
		".":                          {".gitignore", "README.md", "Dockerfile", "docker-compose.yml", "Makefile", ".env"},
	}

	if frontend == "htmx" {
		structure["templates"] = []string{"base.html", "index.html"}
		structure["static"] = []string{
			"css/style.css",
			"js/htmx.min.js",
		}
	}

	err = file.CreateStructure(projectTitle, structure, tmpl, "", map[string]interface{}{
		"Framework": framework,
		"Frontend":  frontend,
	})
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}

	// Initialize Go module first
	err = initGoModule(projectTitle)
	if err != nil {
		log.Fatalf("Failed to initialize go module: %v", err)
	}

	// Create React frontend with Vite if framework is react
	if frontend == "react" {
		fmt.Println("Creating React frontend with Vite...")
		cmd := exec.Command("npm", "create", "vite@latest", "frontend", "--", "--template", "react")
		cmd.Dir = projectTitle
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: Failed to create React frontend: %v\n", err)
			return
		}

		// Install dependencies
		cmd = exec.Command("npm", "install")
		cmd.Dir = filepath.Join(projectTitle, "frontend")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: Failed to install React dependencies: %v\n", err)
			return
		}
	}

	// Initialize Git and add common as submodule synchronously
	repoURL := "https://github.com/indalyadav56/go-common"

	// Remove existing common directory if it exists
	commonDir := filepath.Join(projectTitle, "common")
	if err := os.RemoveAll(commonDir); err != nil {
		log.Fatalf("Error removing existing common directory: %v", err)
	}

	// Initialize git repository
	initCmd := exec.Command("git", "init")
	initCmd.Dir = projectTitle
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr
	if err := initCmd.Run(); err != nil {
		log.Fatalf("Error initializing git repository: %v", err)
	}

	// Add submodule
	cmd := exec.Command("git", "submodule", "add", repoURL, "common")
	cmd.Dir = projectTitle
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error adding common submodule: %v", err)
	}

	// Initialize submodule
	initSubCmd := exec.Command("git", "submodule", "update", "--init", "--recursive")
	initSubCmd.Dir = projectTitle
	initSubCmd.Stdout = os.Stdout
	initSubCmd.Stderr = os.Stderr
	if err := initSubCmd.Run(); err != nil {
		log.Fatalf("Error initializing submodule: %v", err)
	}

	// Run go mod tidy after submodule is set up
	err = runGoModTidy(projectTitle)
	if err != nil {
		log.Fatalf("Failed to run 'go mod tidy': %v", err)
	}
}

func runGoModTidy(basePath string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = basePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running 'go mod tidy': %w", err)
	}

	return nil
}

func initGoModule(projectTitle string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}
	fmt.Println("projectTitle", projectTitle)

	customDir := filepath.Join(currentDir, projectTitle)
	fmt.Println("customDir", customDir)

	// Initialize go.mod file
	cmd := exec.Command("go", "mod", "init", projectTitle)
	cmd.Dir = customDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Try to clean up if initialization fails
		os.RemoveAll(customDir)
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	// Add common module replacement and requirement to main go.mod
	goModContent := fmt.Sprintf(`module %s
	go 1.23.1

	replace common => ./common

	require (
		common v1.0.0
	)
	`, projectTitle)

	err = os.WriteFile(filepath.Join(customDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing go.mod file: %v", err)
	}

	return nil
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
