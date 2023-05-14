# syntax = docker/dockerfile:1

##### PRE BUILD ARGUMENTS ######
ARG UID=1001
ARG GID=1001
ARG USER=nonroot
ARG BINARY="polygon-client"
ARG GOLANG_VERSION=1.20.4
ARG ALPINE_VERSION=3.18
############### BASE IMAGE ################
FROM --platform=$BUILDPLATFORM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS base
ARG TARGETOS 
ARG TARGETARCH
ARG UID
ARG GID
ARG USER
ARG BINARY
## Main Packages
RUN apk add --update make protoc protobuf protobuf-dev git build-base bash curl shadow

RUN addgroup -g $GID -S ${USER} && \
    adduser -u $UID -S ${USER} -G ${USER}
# Copy everything of current project

WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

############### BUILD IMAGE ################
FROM base AS build
ARG TARGETOS
ARG TARGETARCH
ARG BUILD_REVISION

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-X main.Build=${BUILD_REVISION} -s -w" -o /out/${BINARY} .

############## RUNTIME IMAGE ###############
FROM --platform=$BUILDPLATFORM scratch
ARG UID
ARG GID
ARG USER
ARG BINARY
WORKDIR /usr/local/bin
## Since scratch containers don't have shell, we need to copy the /etc/passwd
# https://medium.com/@lizrice/non-privileged-containers-based-on-the-scratch-image-a80105d6d341
COPY --from=base /etc/passwd /etc/passwd 
# Copy the binary from the previous stage while giving permissions to the current user.
COPY --chown=${UID}:${GID} --from=build /out/${BINARY} app
# Copy SSL Certificates
COPY --from=base /etc/ssl/certs/ /etc/ssl/certs/
USER ${USER}

################ ENTRYPOINT ##################
ENTRYPOINT ["/usr/local/bin/app"]

