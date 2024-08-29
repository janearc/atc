.ONESHELL:

secrets:
	yq . ${HOME}/.atc/secrets.yml > config/secrets.yml
	@echo "local secrets file is now tainted, use \"make rmsecrets\" to remove before committing"

rmsecrets:
	@echo "removing local secrets"
	cp config/secrets_example.yml config/secrets.yml

docker:
	docker buildx build --no-cache --tag atc:latest --load .

version:
	@echo "Updating version data"
	@echo "version:" > config/version.yml
	@echo "  build_date: \"`date`\"" >> config/version.yml
	@echo "  build: \"`git describe --tags --always`\"" >> config/version.yml
	@echo "  branch: \"`git branch | grep '^*' | cut -d' ' -f 2`\"" >> config/version.yml

build: version secrets docker 
