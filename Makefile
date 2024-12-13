run:
	go run main.go

format:
	go fmt
test-run:
	go test -v ./...

sample: 
	go run main.go init backend --app=auth --app=user --app=todo

# govulncheck ./...

build:
	go build -o ./bin/go-generator
