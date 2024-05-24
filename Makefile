.PHONY: default

default:
	rm -fv build/clockify-export-arm64*
	GOOS=darwin GOARCH=arm64 go build -o build/clockify-export-arm64-v1.0.1 main.go
	cd build && zip clockify-export-arm64-v1.0.1.zip clockify-export-arm64-v1.0.1
