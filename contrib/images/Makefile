all: casino-env

casino-env:
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --tag cosmicbet/casino-env casino-env

.PHONY: all casino-env
