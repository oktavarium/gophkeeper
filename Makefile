VERSION := $$(git describe --tags --abbrev=0)
DATE := $$(date)
build:
	go build -o client -ldflags="-X 'github.com/oktavarium/gophkeeper/internal/shared/buildinfo.Version=${VERSION}' -X 'github.com/oktavarium/gophkeeper/internal/shared/buildinfo.BuildDate=${DATE}'" cmd/client/main.go
	go build -o server -ldflags="-X 'github.com/oktavarium/gophkeeper/internal/shared/buildinfo.Version=${VERSION}' -X 'github.com/oktavarium/gophkeeper/internal/shared/buildinfo.BuildDate=${DATE}'" cmd/server/main.go

gen:
	go generate ./...

run:
	docker-compose up --build --abort-on-container-exit

lint:
	golangci-lint run