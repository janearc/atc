secrets:
	yq . ${HOME}/.atc/secrets.yml > config/secrets.yml
	@echo "local secrets file is now tainted, use \"make rmsecrets\" to remove before committing"

rmsecrets:
	@echo "removing local secrets"
	cp config/secrets_example.yml config/secrets.yml

docker:
	docker buildx build --no-cache --tag atc:latest --load .

build: secrets docker
