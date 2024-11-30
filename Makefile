# Set the default shell to `/bin/sh` for executing commands in the Makefile.
# `/bin/sh` is used as it is lightweight and widely available across UNIX systems.
SHELL = /bin/sh

# Define the project name for easy reference throughout the Makefile.
# This helps in maintaining a consistent project name and avoiding hardcoding it in multiple places.
PROJECT = "hq-go-retrier"

# --- Prepare | Setup -------------------------------------------------------------------------------

.PHONY: prepare
prepare:
	@# Install the latest version of Lefthook (a Git hooks manager) and set it up.
	go install github.com/evilmartians/lefthook@latest && lefthook install

# --- Go(Golang) ------------------------------------------------------------------------------------

# Define common Go commands with variables for reusability and easier updates.
GOCMD=go
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOTEST=$(GOCMD) test

# Define Go build flags for verbosity and linking.
# Verbose flag for Go commands, helpful for debugging and understanding output.
GOFLAGS := -v
# Linker flags:
# - `-s` removes the symbol table for a smaller binary size.
# - `-w` removes DWARF debugging information.
LDFLAGS := -s -w

# Enable static linking on non-macOS platforms.
# This embeds all dependencies directly into the binary, making it more portable.
ifneq ($(shell go env GOOS),darwin)
	LDFLAGS := -extldflags "-static"
endif

# Define Golangci-lint command for linting Go code.
GOLANGCILINTCMD=golangci-lint
GOLANGCILINTRUN=$(GOLANGCILINTCMD) run

# --- Go Module Management

# Tidy Go modules
# This cleans up `go.mod` and `go.sum` by removing unused dependencies
# and ensuring that only the required packages are listed.
.PHONY: go-mod-tidy
go-mod-tidy:
	$(GOMOD) tidy

# Update Go modules
# Updates all Go dependencies to their latest versions, including both direct and indirect dependencies.
# Useful for staying up-to-date with upstream changes and bug fixes.
.PHONY: go-mod-update
go-mod-update:
	@# Update test dependencies.
	$(GOGET) -f -t -u ./...
	@# Update all other dependencies.
	$(GOGET) -f -u ./...

# --- Go Code Quality and Testing

# Format Go code
# Formats all Go source files in the current module according to Go's standard rules.
# Consistent formatting is crucial for code readability and collaboration.
.PHONY: go-fmt
go-fmt:
	$(GOFMT) ./...

# Lint Go code
# Runs static analysis checks on the Go code using Golangci-lint.
# Ensures the code adheres to best practices and is free from common issues.
# This target also runs `go-fmt` beforehand to ensure the code is formatted.
.PHONY: go-lint
go-lint: go-fmt
	$(GOLANGCILINTRUN) $(GOLANGCILINT) ./...

# Run Go tests
# Executes all unit tests in the module with detailed output.
# The `GOFLAGS` variable is used to enable verbosity, making it easier to debug test results.
.PHONY: go-test
go-test:
	$(GOTEST) $(GOFLAGS) ./...

# --- Help -----------------------------------------------------------------------------------------

# Display help information
# This target prints out a detailed list of all available Makefile commands for ease of use.
# It's a helpful reference for developers using the Makefile.
.PHONY: help
help:
	@echo ""
	@echo "*****************************************************************************"
	@echo ""
	@echo "PROJECT : $(PROJECT)"
	@echo ""
	@echo "*****************************************************************************"
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo " Preparation | Setup:"
	@echo "  prepare .................. prepare repository."
	@echo ""
	@echo " Go Commands:"
	@echo "  go-mod-tidy .............. Tidy Go modules."
	@echo "  go-mod-update ............ Update Go modules."
	@echo "  go-fmt ................... Format Go code."
	@echo "  go-lint .................. Lint Go code."
	@echo "  go-test .................. Run Go tests."
	@echo ""
	@echo " Help Commands:"
	@echo "  help ..................... Display this help information"
	@echo ""

# Set the default goal to `help`.
# This ensures that running `make` without arguments will display the help information.
.DEFAULT_GOAL = help