# ==============================================================================
# Modules support
tidy:
	@echo Running go mod tidy...
	@go mod tidy

# ==============================================================================
# Build commands
tracking-win: tidy
	@echo Building for windows...
	@GOOS=windows GOARCH=386 go build -o grtracking.exe ./tracking

tracking-mac: tidy
	@echo Building for mac...
	@GOOS=darwin GOARCH=amd64 go build -o grtracking ./tracking

tracking-linux: tidy
	@echo Building for linux...
	@GOOS=linux GOARCH=amd64 go build -o grtracking ./tracking