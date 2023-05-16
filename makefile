APP_NAME:=sarif-to-markdown-table

# build for linux amd64
build-linux:
	@echo "Building for linux amd64"
	@GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux-amd64 .

# build for arm mac
build-mac:
	@echo "Building for mac arm64"
	@GOOS=darwin GOARCH=arm64 go build -o bin/$(APP_NAME)-mac-arm64 .