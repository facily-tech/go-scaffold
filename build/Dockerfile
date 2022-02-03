ARG VARIANT="1.17"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT} as development

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10
ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>
RUN apt-get update && apt-get -y install netcat

# hot reload for go https://github.com/cosmtrek/air
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
# install swag cli
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/cweill/gotests/gotests@latest
RUN go install github.com/golang/mock/mockgen@v1.6.0
RUN go get github.com/josharian/impl
# allow any user to write on /go/pkg
RUN chmod -R a+w /go/pkg

# Prod

## Builder binary for Prod
FROM golang:1.17-alpine AS builder

WORKDIR /api
RUN apk add --no-cache git make

COPY go.mod go.sum /
RUN go mod download

COPY . .
RUN make build

## Runner image for prod

FROM alpine as runner
WORKDIR /app

COPY run.sh .
COPY --from=builder /api/bin/api .
RUN chmod +x /app/run.sh

EXPOSE 8080

ENTRYPOINT /app/run.sh /app/api