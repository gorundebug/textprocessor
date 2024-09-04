.PHONY: all

all: gen

gen:
	@mkdir -p "$(PROJECT_DIR)/services/charsprocessor/$(GENERATED_DIR)/grpcsvc"; \
	find . -maxdepth 1 -name "*.proto" | while read -r file; do \
		protoc -I=. -I="$(PROJECT_DIR)" --go_out="$(PROJECT_DIR)/services/charsprocessor/$(GENERATED_DIR)/grpcsvc" --go_opt=paths=source_relative --go-grpc_out="$(PROJECT_DIR)/services/charsprocessor/$(GENERATED_DIR)/grpcsvc" --go-grpc_opt=paths=source_relative $$file; \
	done;