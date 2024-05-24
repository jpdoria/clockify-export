.PHONY: default

default:
	GOOS=darwin GOARCH=arm64 go build -o build/clockify-export-arm64 main.go
