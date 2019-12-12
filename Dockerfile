FROM node:13-alpine as client
RUN apk --update add git
WORKDIR /client
COPY client/package* ./
RUN npm ci
COPY client /client
RUN npm run build:prod

FROM golang:1.13 as server
WORKDIR /server

COPY server/go.* /server/
RUN go mod download \
    && go get -u github.com/rakyll/statik

COPY server /server
COPY --from=client /client/dist /server/dist

RUN statik --src dist

RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o app

FROM alpine:3.9
WORKDIR /reverse-proxy
RUN apk --update add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk add --update ca-certificates \
  && update-ca-certificates \
  && rm -rf /var/cache/apk/*

COPY --from=server /server/app /reverse-proxy
COPY server/config.yaml /reverse-proxy

EXPOSE 8080
EXPOSE 80
EXPOSE 443

ENTRYPOINT /reverse-proxy/app





