FROM golang:1.13.1-alpine AS build
WORKDIR /usr/src
RUN apk add --no-cache git gcc build-base upx
ARG GOVERSION=1.13
ARG PKG
ARG APPNAME
ARG COMMITHASH
ARG BUILDTIME
ARG VERSION
COPY . ./
RUN go build -o app -ldflags "\
        -X ${PKG}/version.APPNAME=${APPNAME} \
        -X ${PKG}/version.VERSION=${VERSION} \
        -X ${PKG}/version.GOVERSION=${GOVERSION} \
        -X ${PKG}/version.BUILDTIME=${BUILDTIME} \
        -X ${PKG}/version.COMMITHASH=${COMMITHASH} \
        -s -w"
RUN upx --best app

FROM alpine
EXPOSE 3000
WORKDIR /usr/src
COPY --from=build /usr/src/app /usr/src/
CMD ["/usr/src/app"]
