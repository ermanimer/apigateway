FROM golang:latest
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build ./cmd/apigateway

EXPOSE 8080

ENTRYPOINT ["/app/apigateway"]