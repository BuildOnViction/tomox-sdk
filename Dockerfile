# build stage
FROM golang

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o backend

EXPOSE 8080
ENTRYPOINT ["/app/backend"]