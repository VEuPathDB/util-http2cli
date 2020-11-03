VERSION=$(shell git describe --tags 2>/dev/null || echo "snapshot")
FILES=$(shell find . -name '*.go')

.PHONY: release
release: bin/http2cli.darwin.$(VERSION).x64.tar.gz bin/http2cli.linux.$(VERSION).x64.tar.gz bin/http2cli.windows.$(VERSION).x64.zip

.PHONY: linux
linux: bin/http2cli

bin/http2cli: $(FILES)
	CGO_ENABLED=0 GOOS=linux go build -o bin/http2cli -ldflags "-X main.version=${VERSION}" v1/cmd/server/main.go

bin/http2cli-mac: $(FILES)
	CGO_ENABLED=0 GOOS=darwin go build -o bin/http2cli-mac -ldflags "-X main.version=${VERSION}" v1/cmd/server/main.go

bin/http2cli.exe: $(FILES)
	CGO_ENABLED=0 GOOS=windows go build -o bin/http2cli.exe -ldflags "-X main.version=${VERSION}" v1/cmd/server/main.go

bin/http2cli.linux.$(VERSION).x64.tar.gz: bin/http2cli
	cd bin && tar -czf http2cli.linux.$(VERSION).x64.tar.gz http2cli
	rm bin/http2cli

bin/http2cli.darwin.$(VERSION).x64.tar.gz: bin/http2cli-mac
	cd bin && tar -czf http2cli.darwin.$(VERSION).x64.tar.gz http2cli-mac
	rm bin/http2cli-mac

bin/http2cli.windows.$(VERSION).x64.zip: bin/http2cli.exe
	cd bin && zip -9 http2cli.darwin.$(VERSION).x64.zip http2cli.exe
	rm bin/http2cli.exe