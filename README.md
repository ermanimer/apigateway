# apigateway

[![Test](https://github.com/ermanimer/apigateway/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/ermanimer/apigateway/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/ermanimer/apigateway/graph/badge.svg?token=rbFp8CZIRk)](https://codecov.io/gh/ermanimer/apigateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/ermanimer/apigateway)](https://goreportcard.com/report/github.com/ermanimer/apigateway)

Apigateway is an API gateway designed for use in development environments or sandboxes. It follows these principles in its design:

 - Its configuration is simple.
 - It can be used directly within Docker Compose.
 - It does not perform load balancing.
 - It does not modify requests and responses.
 - It leverages the reverse proxy provided in Go's standard library, ensuring optimal performance and efficiency.

 # Using with Docker Compose

**docker-compose.yml**

```yaml
version: '3'

services:
  apigateway:
    container_name: apigateway
    restart: always
    image: imererman/apigateway:latest
    volumes:
      - ./config.yml:/app/config.yml
    networks:
      - apigateway-network
    ports:
      - 80:8080
    
  service1:
    container_name: service1
    ...
    networks:
      - apigateway-network

networks:
  apigateway-network:
    name: apigateway-network
```

**config.yml**

```yaml
upstreams:
  - pattern: /service1/
    strip_prefix: true
    url: http://service1:8080
```

With this configuration in docker-compose, all requests starting with the **/service1/** pattern coming to the API gateway will be forwarded to service1 after stripping the pattern. For example:

```
http://localhost:80/service1/health-check -> http://service1/8080/health-check:
```

If ```strip_prefix``` was ```false```, the request coming to the API gateway as above would be forwarded to service1 without stripping the pattern, as follows:

```
http://localhost:80/service1/health-check -> http://service1/service1/8080/health-check:
```

# Configuration

As seen in the example above, only upstream configuration is sufficient. Below, the default configuration for the server is shown. If desired, the server can also be configured.

```yaml
server:
  address: :8080
  read_timeout: 5s
  write_timeout: 10s
  idle_timeout: 120s
  max_header_bytes: 1048576 # 1 MB
  shutdown_timeout: 10s

upstreams:
...
```

# Build And Run

To build the project with Go, you can run the following commands in the project directory.

**build:**

```bash
go build ./cmd/apigateway
```

**run:**

```bash
./apigateway
```

Note: apigateway always looks for the config.yml file in its directory.

# Contribution

Open an issue and let's collectively decide on the changes and features you want. Then, you can proceed by opening a pull request from the main branch of your forked repository to the main branch of the main repository.

# License

This repo is under the MIT license.