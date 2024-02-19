build:
	go build -o client cmd/client/main.go
	go build -o server cmd/server/main.go

gen:
	go generate ./...