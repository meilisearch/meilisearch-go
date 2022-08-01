FROM golang:1.17.11-buster

WORKDIR /home/package

COPY go.mod .
COPY go.sum .

COPY --from=golangci/golangci-lint:v1.42.0 /usr/bin/golangci-lint /usr/local/bin/golangci-lint

RUN go mod download
RUN go mod verify
