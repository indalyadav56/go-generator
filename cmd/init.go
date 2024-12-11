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
	fmt.Println("Creating project")

	// Debug: List all files in embedded FS
	entries, err := templates.TemplateFS.ReadDir(".")
	if err != nil {
		panic(err)
	}
	fmt.Println("Available templates:")
	for _, entry := range entries {
		fmt.Printf("- %s\n", entry.Name())
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
		"templates/gin/main.tmpl",
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
		"config":                            {"env.go", "app.go", "router.go"},
		"db/postgres":                       {"postgres.go", "db_logger.go"},
		"db":                                {"postgres.go"},
		"migrations":                        {""},
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
	fmt.Println("projectTitle", projectTitle)

	fmt.Println("currentDir", currentDir)

	customDir := filepath.Join(currentDir, projectTitle)
	fmt.Println("customDir", customDir)

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
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	// Get the absolute path to the project root directory
	projectRoot := filepath.Dir(filepath.Dir(currentDir))
	srcDir := filepath.Join(projectRoot, "pkg")

	err = CopyFolder(srcDir, projectTitle+"/pkg")
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
