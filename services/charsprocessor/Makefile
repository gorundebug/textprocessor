BINARY_NAME=$(SERVICE_NAME)

.PHONY: all build clean run

all: build

build: gen service_build

service_build:
	@echo "Building $(SERVICE_NAME)..."
	go build -o "$(BIN_DIR)/$(SERVICE_NAME)" "../../cmd/$(SERVICE_NAME)/main.go"

clean:
	@echo "Cleaning $(SERVICE_NAME)..."
	rm -rf "$(BIN_DIR)/$(BINARY_NAME)"

gen:
	@if [ -d proto ]; then \
		echo "Generate proto code..."; \
		find proto -type f -name 'Makefile' | while read makefile; do \
			dir=$$(dirname $$makefile); \
			echo "Generating files in $$dir..."; \
			$(MAKE) -C $$dir gen PROJECT_DIR="$(PROJECT_DIR)" GENERATED_DIR="$(GENERATED_DIR)"; \
		done \
	fi

	@if [ -d grpc ]; then \
		echo "Generate grpc code..."; \
		find grpc -type f -name 'Makefile' | while read makefile; do \
			dir=$$(dirname $$makefile); \
			echo "Generating files in $$dir..."; \
			$(MAKE) -C $$dir gen PROJECT_DIR="$(PROJECT_DIR)" GENERATED_DIR="$(GENERATED_DIR)"; \
		done \
	fi

run:
	go run "../../cmd/$(SERVICE_NAME)/main.go" -config "../../configs/$(SERVICE_NAME)/config.yaml" -values "../../configs/$(SERVICE_NAME)/values.yaml" || true