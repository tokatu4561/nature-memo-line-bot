.PHONY: build

build:
	sam build

start-api:
	sam local start-api --env-vars env.json
	
invoke:
	sam local invoke --env-vars env.json