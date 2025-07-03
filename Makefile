.PHONY:
.SILENT:

build: 
	go build -o ./.bin/auth-service ./cmd/app/main.go

run: build 
	./.bin/auth-service