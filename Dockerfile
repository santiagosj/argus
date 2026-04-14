FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o argus ./cmd/argus

FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite-libs dnsmasq curl nmap

WORKDIR /app
COPY --from=builder /app/argus .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/skills ./skills

EXPOSE 8080
ENV HOME=/app

ENTRYPOINT ["./argus"]
CMD ["demo"]
