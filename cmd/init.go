package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/indalyadav56/go-generator/file"
	"github.com/indalyadav56/go-generator/templates"
	"github.com/spf13/cobra"
)

type ProjectOpts struct {
	Websocket bool
}

var initCmd = &cobra.Command{
	Use:   "new",
	Short: "This command initializes a golang project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		apps, _ := cmd.Flags().GetStringSlice("app")
		framework, _ := cmd.Flags().GetString("framework")
		frontend, _ := cmd.Flags().GetString("frontend")
		websocket, _ := cmd.Flags().GetBool("websocket")
		if len(args) == 0 {
			fmt.Println("provide the name of the project")
			return
		}
		name := args[0]

		CreateProject(strings.ToLower(name), framework, frontend, apps, ProjectOpts{Websocket: websocket})

		if len(apps) > 0 {
			// create apps
			for _, v := range apps {
				CreateApp(strings.ToLower(v), fmt.Sprintf("./%s/internal", name))
			}

			go initSwagger(strings.ToLower(name))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringSliceP("app", "", []string{}, "app names to initialize")
	initCmd.Flags().StringP("framework", "", "", "web framework to use (e.g., htmx)")
	initCmd.Flags().StringP("frontend", "", "", "frontend framework to use (e.g., react, htmx)")
	initCmd.Flags().BoolP("websocket", "w", false, "enable WebSocket support")

}

func CreateProject(projectTitle string, framework string, frontend string, apps []string, opts ProjectOpts) {
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
		"templates/migrations/users.sql.tmpl",
		"templates/nginx/nginx-conf.tmpl",
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

	rootFiles := []string{
		".gitignore",
		"README.md",
		"Dockerfile",
		"docker-compose.yml",
		"Makefile",
		".env",
	}

	for _, app := range apps {
		rootFiles = append(rootFiles, fmt.Sprintf("%s.http", app))
	}

	structure := file.DirectoryStructure{
		fmt.Sprintf("cmd/%s", "api"): {"main.go"},
		"config":                     {"config.go"},
		"docs":                       {""},
		"scripts":                    {"build.sh"},
		"logs":                       {""},
		"migrations":                 []string{""},
		"internal/app":               {"app.go", "deps.go"},
		".":                          rootFiles,
		"nginx":                      {"nginx.conf"},
	}

	if frontend == "htmx" {
		structure["web/templates"] = []string{"base.html", "index.html"}
		structure["web/static"] = []string{
			"css/style.css",
			"js/htmx.min.js",
		}
	}

	if opts.Websocket && framework == "gin" {
		structure["internal/websocket"] = []string{"client.go", "config.go", "server.go", "routes.go", "hub.go", "handler.go"}
	}

	initialApps := make(map[string]bool)
	for _, app := range apps {
		if app == "user" {
			initialApps["user"] = true
		}
		if app == "auth" {
			initialApps["auth"] = true
		}
	}

	err = file.CreateStructure(projectTitle, structure, tmpl, "", map[string]interface{}{
		"Framework":   framework,
		"Frontend":    frontend,
		"InitialApps": initialApps,
	})
	if err != nil {
		log.Fatalf("Failed to create structure: %v\n", err)
	}

	go func() {
		for _, app := range apps {
			if app != "auth" && app != "authentication" {
				if app == "user" {
					migrationName := app
					cmd := exec.Command("goose", "create", migrationName, "sql")
					cmd.Dir = filepath.Join(projectTitle, "migrations")
					if err := cmd.Run(); err != nil {
						log.Printf("Failed to create migration for %s: %v", app, err)
						continue
					}
					time.Sleep(1 * time.Second)
					output := new(bytes.Buffer)

					files, err := filepath.Glob(filepath.Join(projectTitle, "migrations", fmt.Sprintf("*_%s.sql", migrationName)))
					if err != nil || len(files) == 0 {
						log.Printf("Failed to find migration file for %s: %v", app, err)
						continue
					}

					for _, file := range files {
						if strings.Contains(file, "user.sql") {
							if err := tmpl.ExecuteTemplate(output, "sql_migration", map[string]string{"AppName": app}); err != nil {
								log.Printf("Failed to generate migration content for %s: %v", app, err)
								continue
							}

							if err := os.WriteFile(file, output.Bytes(), 0644); err != nil {
								log.Printf("Failed to write migration file for %s: %v", app, err)
							}
						}
					}

					output.Reset()
					break
				}

			}
		}
	}()

	// Initialize Go module first
	err = initGoModule(projectTitle)
	if err != nil {
		log.Fatalf("Failed to initialize go module: %v", err)
	}

	// Create React frontend with Vite if framework is react
	if frontend == "react" {
		fmt.Println("Creating React frontend with Vite...")
		cmd := exec.Command("npm", "create", "vite@latest", "frontend", "--", "--template", "react-ts")
		cmd.Dir = projectTitle
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: Failed to create React frontend: %v\n", err)
			return
		}

		// // Install dependencies
		// cmd = exec.Command("npm", "install")
		// cmd.Dir = filepath.Join(projectTitle, "frontend")
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// if err := cmd.Run(); err != nil {
		// 	log.Printf("Warning: Failed to install React dependencies: %v\n", err)
		// 	return
		// }
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

	// Initialize Swagger
	err = initSwagger(projectTitle)
	if err != nil {
		log.Fatalf("Failed to initialize Swagger: %v", err)
	}
}

// initSwagger initializes Swagger documentation for the project
func initSwagger(projectPath string) error {
	// Check if swag CLI tool is installed
	checkCmd := exec.Command("swag", "version")
	if err := checkCmd.Run(); err != nil {
		// Install swag CLI tool if not found
		installCmd := exec.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
		installCmd.Dir = projectPath
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install swag: %w", err)
		}
	}

	// Add Swagger dependencies to go.mod
	swaggerDeps := []string{
		"github.com/swaggo/swag@v1.16.2",
		"github.com/swaggo/gin-swagger@v1.6.0",
		"github.com/swaggo/files@v1.0.1",
	}
	for _, dep := range swaggerDeps {
		checkCmd := exec.Command("go", "list", "-m", dep)
		checkCmd.Dir = projectPath
		if err := checkCmd.Run(); err != nil {
			// Dependency not found, install it
			getCmd := exec.Command("go", "get", dep)
			getCmd.Dir = projectPath
			if err := getCmd.Run(); err != nil {
				return fmt.Errorf("failed to add swagger dependency %s: %w", dep, err)
			}
		}
	}

	// Run swag init
	cmd := exec.Command("swag", "init", "-g", "cmd/api/main.go", "-o", "./docs")
	// cmd := exec.Command("swag", "init", "-g", "cmd/api/main.go")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize swagger: %w", err)
	}

	return nil
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
