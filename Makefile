run:
	go run main.go

format:
	go fmt
test-run:
	go test -v ./...

# go run main.go init backend --app=auth --app=user --app=todo

govulncheck ./...