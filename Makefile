.PHONY: default

VERSION := 0.4.5

default:
	rm -fv bin/clockify-export-*
	GOOS=darwin GOARCH=arm64 go build -o bin/clockify-export-arm64-$(VERSION) -ldflags "-X main.ver=$(VERSION) -X 'main.build=`date +%Y%m%d%H%M%S%3N`'" main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/clockify-export-amd64-$(VERSION) -ldflags "-X main.ver=$(VERSION) -X 'main.build=`date +%Y%m%d%H%M%S%3N`'" main.go
