# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o backend

# final stage
#FROM golang:1.12.6-alpine
#COPY --from=builder /app/backend /app/
EXPOSE 8080
ENTRYPOINT ["/app/backend"]