SERVICES := wordsprocessor charsprocessor
BIN_DIR := $(abspath bin)
PROJECT_DIR := $(abspath .)
GENERATED_DIR := generated

.PHONY: all build clean run gen $(SERVICES)

all: build

rebuild: clean build

build: printversion gen prepare $(SERVICES)

printversion:
	@go version;

clean: $(addprefix clean-,$(SERVICES))
	echo "Clean project..."
	rm -rf $(BIN_DIR)

run: gen prepare
	@trap 'kill $(jobs -p)' SIGINT; \
	for service in $(SERVICES); do \
		$(MAKE) run-$$service & \
	done; \
	wait

$(SERVICES):
	@bash ./version_inc.sh "./cmd/$@/main.go"
	mkdir -p $(BIN_DIR)/$@
	$(MAKE) -C services/$@ build SERVICE_NAME=$@ BIN_DIR="$(BIN_DIR)/$@" PROJECT_DIR="$(PROJECT_DIR)" GENERATED_DIR="$(GENERATED_DIR)"

prepare:
	@if [ -f ./go.mod ]; then \
		echo "Running go mod tidy..."; \
		go mod tidy; \
	fi

gen: $(addprefix gen-,$(SERVICES))
	@if [ -d pkg/proto ]; then \
		echo "Generate files..."; \
		find pkg/proto -type f -name 'Makefile' | while read makefile; do \
			dir=$$(dirname $$makefile); \
			echo "Generating files in $$dir..."; \
			$(MAKE) -C $$dir gen PROJECT_DIR="$(PROJECT_DIR)" GENERATED_DIR="$(GENERATED_DIR)"; \
		done \
    fi

clean-%:
	$(MAKE) -C services/$* clean SERVICE_NAME=$* BIN_DIR="$(BIN_DIR)" PROJECT_DIR="$(PROJECT_DIR)" GENERATED_DIR="$(GENERATED_DIR)"

run-%: gen prepare
	$(MAKE) -C services/$* run SERVICE_NAME=$* BIN_DIR="$(BIN_DIR)" PROJECT_DIR="$(PROJECT_DIR)" GENERATED_DIR="$(GENERATED_DIR)"

gen-%:
	$(MAKE) -C services/$* gen SERVICE_NAME=$* BIN_DIR="$(BIN_DIR)" PROJECT_DIR="$(PROJECT_DIR)" GENERATED_DIR="$(GENERATED_DIR)"