# --- Build stage ---
FROM golang:1.23.3-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/kirolink kirolink.go

# --- Runtime stage ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates sqlite && \
    adduser -D -u 1000 kirolink

WORKDIR /home/kirolink

COPY --from=builder /out/kirolink /usr/local/bin/kirolink

USER kirolink

EXPOSE 8080

ENTRYPOINT ["kirolink"]
CMD ["server", "8080"]
