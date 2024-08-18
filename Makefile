#-include .env

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
all, help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nMakefile help:\n  make \033[36m<target>\033[0m\n"} /^[0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

generate: generate_grpc generate_ts_client ### generate all
	echo "generate..."
.PHONY: generate

generate_grpc: ### generate grpc
	@-GOBIN=$(GOBIN) go install github.com/bufbuild/buf/cmd/buf@v1.35.1
	cd ./proto && \
	buf dep update && \
	buf generate
	sed -i '' -e 's/https:\/\/localhost/http:\/\/localhost/g' generated/openapiv3/openapi.yaml
.PHONY: generate_grpc

generate_ts_client:
	openapi-generator-cli generate -g typescript-angular -i generated/openapiv3/openapi.yaml -o frontend/src/app/api
.PHONY: .generate_ts_client

deps: deps_brew deps_npm ### install dependencies
	echo "deps..."
.PHONY: deps

deps_brew: ### install dependencies from brew
	brew install bufbuild/buf/buf protobuf openjdk
.PHONY: deps_brew

deps_npm: ### install dependencies from npm
	npm install -g @angular/cli@18
	npm install -g @openapitools/openapi-generator-cli
.PHONY: deps_npm