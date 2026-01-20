FROM golang:1.21-bookworm

WORKDIR /home/package

RUN git config --global --add safe.directory /home/package

COPY go.mod .
COPY go.sum .

COPY --from=golangci/golangci-lint:v2.3.0 /usr/bin/golangci-lint /usr/local/bin/golangci-lint

RUN go mod download
RUN go mod verify
