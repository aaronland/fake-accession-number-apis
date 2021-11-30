cli:
	go build -mod vendor -o bin/import cmd/import/main.go
	go build -mod vendor -o bin/lookup cmd/lookup/main.go
	go build -mod vendor -o bin/server cmd/server/main.go
