FROM golang:latest as builder

RUN apt update && apt upgrade -y
RUN apt install tzdata ca-certificates -y
RUN update-ca-certificates

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/apigateway

FROM alpine:latest as runner

WORKDIR /app

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/apigateway /app/apigateway

EXPOSE 8080

ENTRYPOINT ["/app/apigateway"]