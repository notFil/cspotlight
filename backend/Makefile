APP_EXECUTABLE=app

include .env
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=$(subst ",,$(DB_URL))

.PHONY: serve build run clean up down

serve:
	go run cmd/api/main.go

build:
	GOARCH=amd64 GOOS=darwin go build -o ${APP_EXECUTABLE}-darwin cmd/api/main.go
	GOARCH=amd64 GOOS=linux go build -o ${APP_EXECUTABLE}-linux cmd/api/main.go
	GOARCH=amd64 GOOS=windows go build -o ${APP_EXECUTABLE}-windows cmd/api/main.go

run: build
	./${APP_EXECUTABLE}

clean:
	go clean
	rm ${APP_EXECUTABLE}-darwin
	rm ${APP_EXECUTABLE}-linux
	rm ${APP_EXECUTABLE}-windows

up:
	goose -dir migrations up

down:
	goose -dir migrations down