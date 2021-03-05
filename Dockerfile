FROM golang:1.15.7-alpine AS build
ENV GOVERSION=1.15.7
WORKDIR /usr/src
RUN apk add --no-cache git=2.30.1-r0 gcc=10.2.1_pre1-r3 build-base=0.5-r2 upx=3.96-r0
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

FROM alpine:3
EXPOSE 3000
WORKDIR /usr/src
COPY --from=build /usr/src/app /usr/src/
CMD ["/usr/src/app"]
