# build stage
FROM golang:alpine AS build-env
RUN apk add --no-cache git
RUN apk add --no-cache g++
ADD . /src
RUN cd /src && go get -d -v ./... && go build -o goapp

# final stage
FROM alpine
ENV CONFIG_PATH /app/config.toml
ENV SERVER_PORT 3000
WORKDIR /app
COPY --from=build-env /src/goapp /app/
COPY --from=build-env /src/index.html /app/
COPY --from=build-env /src/favicon.ico /app/
COPY --from=build-env /src/config.toml /app/
ENTRYPOINT ./goapp
