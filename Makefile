run:
	go run main.go

format:
	go fmt
test-run:
	go test -v ./...

sample: 
	go run main.go init backend --app=auth --app=user --app=todo

# ╰─ go run main.go init backend --app=auth --app=notification --app=payment --app=storage --frontend=htmx --framework=gin --driver=postgres --orm=gorm
# go run main.go init my_project --app=auth --app=notification --app=user --framework=gin --frontend=htmx
# govulncheck ./...

build:
	go build -o ./bin/go-generator
