# --- Global -------------------------------------------------------------------
O = out

all: backend frontend ## build test and lint frontend and backend
	@if [ -e .git/rebase-merge ]; then git --no-pager log -1 --pretty='%h %s'; fi
	@echo '$(COLOUR_GREEN)Success$(COLOUR_NORMAL)'

clean::  ## Remove generated files
	-rm -rf $(O)

.PHONY: all clean

# --- Build --------------------------------------------------------------------
# Build all subdirs of ./cmd, excluding those with a leading underscore.
BIN = $(O)/covid19
STUB = frontend/src/data.js

backend: build test check-coverage lint  ## build, test, check coverage, lint backend

build: | $(O)  ## Build binary
	go build -o $(BIN) ./backend

run: build  ## Run binary with testdata
	$(BIN) < backend/testdata/data.csv

data: build  ## Fetch covid-19 CSV data and parse it into api stub
	printf "export const data = " > $(STUB)
	curl https://raw.githubusercontent.com/datasets/covid-19/master/data/time-series-19-covid-combined.csv | $(BIN) >> $(STUB)

.PHONY: build run

# --- Test ---------------------------------------------------------------------
COVERFILE = $(O)/coverage.txt
COVERAGE = 73.8

test: | $(O)  ## Run tests and generate a coverage file
	go test -coverprofile=$(COVERFILE) ./...

check-coverage: test  ## Check that test coverage meets the required level
	@go tool cover -func=$(COVERFILE) | $(CHECK_COVERAGE) || $(FAIL_COVERAGE)

cover: test  ## Show test coverage in your browser
	go tool cover -html=$(COVERFILE)

CHECK_COVERAGE = awk -F '[ \t%]+' '/^total:/ {print; if ($$3 < $(COVERAGE)) exit 1}'
FAIL_COVERAGE = { echo '$(COLOUR_RED)FAIL - Coverage below $(COVERAGE)%$(COLOUR_NORMAL)'; exit 1; }

.PHONY: test check-coverage cover

# --- Lint ---------------------------------------------------------------------
GOLINT_VERSION = 1.30.0
GOLINT_INSTALLED_VERSION = $(or $(word 4,$(shell golangci-lint --version 2>/dev/null)),0.0.0)
GOLINT_MIN_VERSION = $(shell printf '%s\n' $(GOLINT_VERSION) $(GOLINT_INSTALLED_VERSION) | sort -V | head -n 1)
GOPATH1 = $(firstword $(subst :, ,$(GOPATH)))
LINT_TARGET = $(if $(filter $(GOLINT_MIN_VERSION),$(GOLINT_VERSION)),lint-with-local,lint-with-docker)

lint: $(LINT_TARGET)  ## Lint source code

lint-with-local:  ## Lint source code with locally installed golangci-lint
	golangci-lint run

lint-with-docker:  ## Lint source code with docker image of golangci-lint
	docker run --rm -w /src \
		-v $(shell pwd):/src -v $(GOPATH1):/go -v $(HOME)/.cache:/root/.cache \
		golangci/golangci-lint:v$(GOLINT_VERSION) \
		golangci-lint run

.PHONY: lint lint-with-local lint-with-docker

# --- Frontend ----------------------------------------------------------------
frontend: frontend-lint frontend-build  ## Lint and build frontend

serve: frontend-build ## Build app and serve
	yarn serve

dev: frontend-init  ## Start frontend development server
	yarn dev

frontend/node_modules:
	yarn install

frontend-init: frontend/node_modules

frontend-lint: frontend-init  ## Lint frontend
	yarn lint

frontend-build: frontend-init  ## Build frontend
	yarn build

clean::
	rm -rf frontend/node_modules frontend/public/build

.PHONY: dev frontend-build frontend-init frontend-lint frontend serve

# --- Utilities ----------------------------------------------------------------
COLOUR_NORMAL = $(shell tput sgr0 2>/dev/null)
COLOUR_RED    = $(shell tput setaf 1 2>/dev/null)
COLOUR_GREEN  = $(shell tput setaf 2 2>/dev/null)
COLOUR_WHITE  = $(shell tput setaf 7 2>/dev/null)

help:
	@awk -F ':.*## ' 'NF == 2 && $$1 ~ /^[A-Za-z0-9_-]+$$/ { printf "$(COLOUR_WHITE)%-30s$(COLOUR_NORMAL)%s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

$(O):
	@mkdir -p $@

.PHONY: help
