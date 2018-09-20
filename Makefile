PROJECT_NAME = azure-storage-blob-go
WORK_DIR = /go/src/github.com/Azure/${PROJECT_NAME}

define with_docker
	WORK_DIR=$(WORK_DIR) docker-compose run --rm $(PROJECT_NAME) $(1)
endef

login: setup ## get a shell into the container
	WORK_DIR=$(WORK_DIR) docker-compose run --rm --entrypoint /bin/bash $(PROJECT_NAME)

docker-compose:
	which docker-compose

docker-build: docker-compose
	WORK_DIR=$(WORK_DIR) docker-compose build --force-rm

docker-clean: docker-compose
	WORK_DIR=$(WORK_DIR) docker-compose down

setup: clean docker-build

test: setup ## run go tests
	$(call with_docker,go test -race -short -cover -v ./azblob)

build: setup ## build binaries for the project
	GOOS=linux $(call with_docker,go build ./azblob,-e GOOS)
	GOOS=darwin $(call with_docker,go build ./azblob,-e GOOS)
	GOOS=windows $(call with_docker,go build ./azblob,-e GOOS)

all: setup build

clean: docker-clean ## clean environment and binaries
	rm -rf bin

help: ## display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
