# File Tools — common developer commands.
# The shareable Windows installer is built by CI (see .github/workflows), or
# locally on Windows via scripts/build-windows.ps1.

APP := filetools
PKG := ./cmd/filetools

.PHONY: run test build vet fmt tidy icon clean

run: ## Run the app locally
	go run $(PKG)

test: ## Run all tests
	go test ./...

build: ## Build a native binary into dist/
	mkdir -p dist
	go build -o dist/$(APP) $(PKG)

vet: ## Static checks
	go vet ./...

fmt: ## Format the code
	gofmt -w .

tidy: ## Tidy module dependencies
	go mod tidy

icon: ## Regenerate build/appicon.ico from build/appicon.png
	go run build/mkicon.go

clean: ## Remove build output
	rm -rf dist FileTools.exe build/FileToolsSetup.exe
