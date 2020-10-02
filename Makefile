
# VARIABLES
# -


# ENVIRONMENT VARIABLES
export GO111MODULE=on
export GOPRIVATE=gitlab.com/swissblock


# CONFIG
.PHONY: help print-variables
.DEFAULT_GOAL = help


# ACTIONS

build :		## Build package
	go build ./...

mod-down :		## Download go modules references
	go mod download

mod-tidy :		## Tidy go modules references
	go mod tidy

test:
	go test -coverprofile=coverage.out -count=5 -race ./...

## helpers

available-exchange :		## Print available exchanges
	echo $(AVAILABLE_EXCHANGES)

help :		## Help
	@echo ""
	@echo "*** \033[33mMakefile help\033[0m ***"
	@echo ""
	@echo "Targets list:"
	@grep -E '^[a-zA-Z_-]+ :.*?## .*$$' $(MAKEFILE_LIST) | sort -k 1,1 | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""

print-variables :		## Print variables values
	@echo ""
	@echo "*** \033[33mMakefile variables\033[0m ***"
	@echo ""
	@echo "- - - makefile - - -"
	@echo "MAKE: $(MAKE)"
	@echo "MAKEFILES: $(MAKEFILES)"
	@echo "MAKEFILE_LIST: $(MAKEFILE_LIST)"
	@echo ""
