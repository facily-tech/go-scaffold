FROM golang:1.17-alpine AS builder

WORKDIR /api
RUN apk add --no-cache git make

COPY go.mod go.sum /
RUN go mod download

COPY . .
RUN make build

# runner

FROM alpine as runner
WORKDIR /go/bin

COPY --from=builder /api/bin/api .
EXPOSE 8080

ENTRYPOINT /go/bin/api