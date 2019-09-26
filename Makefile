.PHONY: build

build: build-client statik build-server

.PHONY: build-client
build-client:
	cd client &&\
	yarn build:prod

.PHONY: build-server
build-server:
	cd server &&\
	go build -o app

.PHONY: statik
statik:
	cd server &&\
	statik --src ../client/dist


