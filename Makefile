.PHONY: default

default:
	rm -fv bin/clockify-export-*
	GOOS=darwin GOARCH=arm64 go build -o bin/clockify-export-arm64-0.4.3 -ldflags "-X main.ver=0.4.3 -X 'main.build=`date +%Y%m%d%H%M%S%3N`'" main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/clockify-export-amd64-0.4.3 -ldflags "-X main.ver=0.4.3 -X 'main.build=`date +%Y%m%d%H%M%S%3N`'" main.go
