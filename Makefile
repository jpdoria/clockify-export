.PHONY: default

default:
	rm -fv build/clockify-export-*
	GOOS=darwin GOARCH=arm64 go build -o build/clockify-export-arm64-0.3.0 main.go
	GOOS=darwin GOARCH=amd64 go build -o build/clockify-export-amd64-0.3.0 main.go
