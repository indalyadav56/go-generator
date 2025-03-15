run:
	go run main.go

format:
	go fmt
test-run:
	go test -v ./...

sample: 
	go run main.go new backend --app=auth --app=user --app=todo

# ╰─ go run main.go new backend --app=auth --app=notification --app=payment --app=storage --frontend=htmx --framework=gin --driver=postgres --orm=gorm
# go run main.go new my_project --app=auth --app=notification --app=user --framework=gin --frontend=htmx
# govulncheck ./...
# --websocket

build:
	go build -o ./bin/go-generator

# docker-compose down --volumes
# docker-compose up --build