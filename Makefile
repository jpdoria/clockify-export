.PHONY: default

default:
	rm -fv bin/clockify-export-*
	GOOS=darwin GOARCH=arm64 go build -o bin/clockify-export-arm64-0.4.1 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/clockify-export-amd64-0.4.1 main.go
