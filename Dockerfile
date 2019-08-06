# build stage
FROM golang:1.12.5

LABEL maintainer="Hai Dam <haidv@tomochain.com>"

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o backend

EXPOSE 8080
ENTRYPOINT ["/app/backend"]