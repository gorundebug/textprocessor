.PHONY: all

all: gen

gen:
	@mkdir -p "$(PROJECT_DIR)/services/wordsprocessor/$(GENERATED_DIR)/pb"; \
	find . -maxdepth 1 -name "*.proto" | while read -r file; do \
		protoc -I=. -I="$(PROJECT_DIR)" --go_out="$(PROJECT_DIR)/services/wordsprocessor/$(GENERATED_DIR)/pb" --go_opt=paths=source_relative $$file; \
	done;