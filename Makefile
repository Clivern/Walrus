GO           ?= go
GOFMT        ?= $(GO)fmt
NPM          ?= npm
NPX          ?= npx
RHINO        ?= rhino
pkgs          = ./...
HUGO ?= hugo


help: Makefile
	@echo
	@echo " Choose a command run in Walrus:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo


## install_revive: Install revive for linting.
install_revive:
	@echo ">> ============= Install Revive ============= <<"
	$(GO) get github.com/mgechev/revive


## style: Check code style.
style:
	@echo ">> ============= Checking Code Style ============= <<"
	@fmtRes=$$($(GOFMT) -d $$(find . -path ./vendor -prune -o -name '*.go' -print)); \
	if [ -n "$${fmtRes}" ]; then \
		echo "gofmt checking failed!"; echo "$${fmtRes}"; echo; \
		echo "Please ensure you are using $$($(GO) version) for formatting code."; \
		exit 1; \
	fi


## check_license: Check if license header on all files.
check_license:
	@echo ">> ============= Checking License Header ============= <<"
	@licRes=$$(for file in $$(find . -type f -iname '*.go' ! -path './vendor/*') ; do \
               awk 'NR<=3' $$file | grep -Eq "(Copyright|generated|GENERATED)" || echo $$file; \
       done); \
       if [ -n "$${licRes}" ]; then \
               echo "license header checking failed:"; echo "$${licRes}"; \
               exit 1; \
       fi


## test_short: Run test cases with short flag.
test_short:
	@echo ">> ============= Running Short Tests ============= <<"
	$(GO) test -short $(pkgs)


## test: Run test cases.
test:
	@echo ">> ============= Running All Tests ============= <<"
	$(GO) test -v -cover $(pkgs)


## lint: Lint the code.
lint:
	@echo ">> ============= Lint All Files ============= <<"
	revive -config config.toml -exclude vendor/... -formatter friendly ./...


## verify: Verify dependencies
verify:
	@echo ">> ============= List Dependencies ============= <<"
	$(GO) list -m all
	@echo ">> ============= Verify Dependencies ============= <<"
	$(GO) mod verify


## format: Format the code.
format:
	@echo ">> ============= Formatting Code ============= <<"
	$(GO) fmt $(pkgs)


## vet: Examines source code and reports suspicious constructs.
vet:
	@echo ">> ============= Vetting Code ============= <<"
	$(GO) vet $(pkgs)


## coverage: Create HTML coverage report
coverage:
	@echo ">> ============= Coverage ============= <<"
	rm -f coverage.html cover.out
	$(GO) test -coverprofile=cover.out $(pkgs)
	go tool cover -html=cover.out -o coverage.html


## ci: Run all CI tests.
ci: style check_license test vet lint
	@echo "\n==> All quality checks passed"


## serve_ui: Serve admin dashboard
serve_ui:
	cd web;$(NPM) run serve


## build_ui: Builds admin dashboard for production
build_ui:
	cd web;$(NPM) run build


## check_ui_format: Check dashboard code format
check_ui_format:
	cd web;$(NPX) prettier  --check .


## format_ui: Format dashboard code
format_ui:
	cd web;$(NPX) prettier  --write .


## api_mock: API mock server
api_mock:
	$(RHINO) serve -c mocks/.rhino.json


## run: Run the service
run:
	$(GO) run walrus.go


.PHONY: help
