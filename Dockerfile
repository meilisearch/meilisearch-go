FROM golang:1.20-buster

WORKDIR /home/package

COPY go.mod .
COPY go.sum .

COPY --from=golangci/golangci-lint:v2.3.0 /usr/bin/golangci-lint /usr/local/bin/golangci-lint

RUN go mod download
RUN go mod verify
