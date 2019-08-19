# build stage
FROM golang:1.12-alpine@sha256:1121c345b1489bb5e8a9a65b612c8fed53c175ce72ac1c76cf12bbfc35211310 as builder

ENV GO111MODULE=on

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
RUN apk update && apk add --no-cache git \
                                     gcc \
                                     musl-dev \
                                     linux-headers \
                                     tzdata

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o backend

FROM scratch AS final

LABEL author="Hai Dam <haidv@tomochain.com>"

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /app/backend /

WORKDIR /

USER nobody:nobody

ENTRYPOINT ["/backend"]

EXPOSE 8080