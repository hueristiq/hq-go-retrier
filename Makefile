SHELL = /bin/sh

PROJECT = "hq-go-retrier"

# ------------------------------------------------------------------------------------------------------------------------------
# --- Prepare | Setup ----------------------------------------------------------------------------------------------------------
# ------------------------------------------------------------------------------------------------------------------------------

.PHONY: prepare
prepare:
	@# Install the latest version of Lefthook (a Git hooks manager) and set it up.
	go install github.com/evilmartians/lefthook@latest && lefthook install

# ------------------------------------------------------------------------------------------------------------------------------
# --- Go (Golang) --------------------------------------------------------------------------------------------------------------
# ------------------------------------------------------------------------------------------------------------------------------

GOCMD=go
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOTEST=$(GOCMD) test

GOFLAGS := -v
LDFLAGS := -s -w
ifneq ($(shell go env GOOS),darwin)
	LDFLAGS := -extldflags "-static"
endif

GOLANGCILINTCMD=golangci-lint
GOLANGCILINTRUN=$(GOLANGCILINTCMD) run

.PHONY: go-mod-tidy
go-mod-tidy:
	$(GOMOD) tidy

.PHONY: go-mod-update
go-mod-update:
	@# Update test dependencies.
	$(GOGET) -f -t -u ./...
	@# Update all other dependencies.
	$(GOGET) -f -u ./...

.PHONY: go-fmt
go-fmt:
	$(GOFMT) ./...

.PHONY: go-lint
go-lint: go-fmt
	$(GOLANGCILINTRUN) $(GOLANGCILINT) ./...

.PHONY: go-test
go-test:
	$(GOTEST) $(GOFLAGS) ./...

# ------------------------------------------------------------------------------------------------------------------------------
# --- Help ---------------------------------------------------------------------------------------------------------------------
# ------------------------------------------------------------------------------------------------------------------------------

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

.DEFAULT_GOAL = help