.DEFAULT_GOAL:=help

# make -f spanner/Makefile test-spanner

.PHONY: help
help: ## Show this help.
	@grep -E '^[0-9a-zA-Z_%-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%s\033[0m\n\t%s\n", $$1, $$2}'


.PHONY: spanner-emulator-up
spanner-emulator-up: ## initialize Spanner emulator database for develop. run in service work
	docker compose -f .devcontainer/docker-compose.yaml up -d --build
	docker compose -f .devcontainer/docker-compose.yaml exec work \
		make spanner-emulator-docker


.PHONY: spanner-emulator-exec
spanner-emulator-exec: ## initialize Spanner emulator database for develop. run in service work
	docker compose -f .devcontainer/docker-compose.yaml exec work bash


.PHONY: spanner-emulator-docker
spanner-emulator-docker: ## initialize Spanner emulator database for develop. run in service work
	gcloud config set project spanner-dump-where
	gcloud config set auth/disable_credentials true
	yes | gcloud config set api_endpoint_overrides/spanner http://spanner:9020/
	SPANNER_EMULATOR_HOST=spanner:9010 \
		yes | gcloud spanner instances delete example || true
	SPANNER_EMULATOR_HOST=spanner:9010 \
		gcloud spanner instances create example --config=emulator-config --description="Instance for integration test"
	SPANNER_EMULATOR_HOST=spanner:9010 \
		gcloud spanner databases create db --instance=example
