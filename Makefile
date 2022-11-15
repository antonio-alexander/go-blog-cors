## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

swagger_port = 10000 # defaults

# REFERENCE: https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop
help: ## - Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

# Reference: https://medium.com/@pedram.esmaeeli/generate-swagger-specification-from-go-source-code-648615f7b9d9
check-swagger: ## - validate/install swagger (v0.29.0)
	@which swagger || (go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0)

swagger: check-swagger ## - generate the swagger.json
	@swagger generate spec --work-dir=./internal/swagger --output ./tmp/swagger.json --scan-models

validate-swagger: swagger ## - validate the swagger.json
	@swagger validate ./tmp/swagger.json

serve-swagger: swagger ## - serve (web) the swagger.json
	@swagger serve -F=swagger ./tmp/swagger.json -p ${swagger_port} --no-open

build: ## - build the source (latest)
	@docker compose build --build-arg GIT_COMMIT=`git rev-parse HEAD` --build-arg GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
	@docker image prune -f

run: ## - run only the dependencies (docker) detached
	@docker compose up --detach --wait

stop: ## - stop and remove all the containers in the compose
	@docker compose down
	@docker compose rm -f
