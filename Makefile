# Makefile for the project

GO := go
GOBUILD := $(GO) build
GOFILES := $(shell find . -name "*.go" -type f)
GOLANGCI_LINT_VERSION := v1.55.1
GOLANGCI_LINT_FILE := bin/golangci-lint
GOLANGCI_LINT_VERSIONED := $(GOLANGCI_LINT_FILE)-$(GOLANGCI_LINT_VERSION)
GOLINT := $(GOLANGCI_LINT_VERSIONED) run
GOAIR := bin/air
TAILWIND := npx tailwindcss

# Go mod tidy
.PHONY: tidy
tidy:
	@echo "Tidying up the go.mod and go.sum files..."
	@$(GO) mod tidy

$(GOLANGCI_LINT_VERSIONED):
	@echo "Setting up golangci-lint..."
	@mkdir -p bin
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin $(GOLANGCI_LINT_VERSION)
	@mv $(GOLANGCI_LINT_FILE) $(GOLANGCI_LINT_VERSIONED)

$(GOAIR):
	@echo "Setting up air for hot reloading..."
	@mkdir -p bin
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b ./bin

.PHONY: setup
setup: $(GOLANGCI_LINT_VERSIONED) $(GOAIR) 
	@echo "Installing tools..."

.PHONY: fmt
fmt: setup
	@echo "Formatting the code..."
	@$(GOLINT) --fix

.PHONY: lint
lint: setup 
	@echo "Linting..."
	@$(GOLINT)

.PHONY: build
build: 
	@echo "Building the application..."
	@$(GOBUILD) -o ./app -v

.PHONY: start
start: build
	@echo "Starting hot reloading server..."
	@$(GOAIR)

.PHONY: css
css:
	@echo "Building tailwind css..."
	@$(TAILWIND) -i css/input.css -o css/output.css --minify

.PHONY: css-watch
css-watch:
	@echo "Watching build tailwind css..."
	@$(TAILWIND) -i css/input.css -o css/output.css --watch
