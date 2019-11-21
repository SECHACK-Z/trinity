.PHONY: build

build: build-client statik build-server

.PHONY: build-client
build-client:
	cd client &&\
	yarn build

.PHONY: build-server
build-server:
	cd server &&\
	go build -o app

.PHONY: statik
statik:
	cd server &&\
	statik --src ../client/dist

.PHONY: pubsub
pubsub:
	genny -in=server/pubsub/pubsub.go -out=server/pubsub/gen-pubsub.go gen "EventType=$(shell awk '/struct/ {print $$2}' server/pubsub/events.go | paste -s -d , -)"


