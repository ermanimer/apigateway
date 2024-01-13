# apigateway

[![Test](https://github.com/ermanimer/apigateway/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/ermanimer/apigateway/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/ermanimer/apigateway/graph/badge.svg?token=rbFp8CZIRk)](https://codecov.io/gh/ermanimer/apigateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/ermanimer/apigateway)](https://goreportcard.com/report/github.com/ermanimer/apigateway)


apigateway is a straightforward API gateway developed in Go.

# Configuration

The application's configuration is outlined in a YAML file, typically named `config.yaml`. This file comprises two primary sections: `server` and `upstreams`.

## Server

The `server` section specifies the server's configuration. All fields are optional, and if not provided, the default values (as shown in the example below) will be used.

```yaml
server:
  address: :8080             # The server will listen on this address
  read_timeout: 5s           # The maximum duration for reading the entire request, including the body
  write_timeout: 10s         # The maximum duration before timing out writes of the response
  idle_timeout: 120s         # The maximum amount of time to wait for the next request when keep-alives are enabled
  max_header_bytes: 1048576  # The maximum size of request headers in bytes
  shutdown_timeout: 10s      # The maximum duration before shutting down the server
```

## Upstreams

The `upstreams` section outlines a list of upstream servers. Each upstream server is defined by a `pattern`, `strip_prefix`, and `url`.

```yaml
upstreams:
  - pattern: /service1/       # The pattern to match in the incoming request
    strip_prefix: true        # If true, the pattern will be removed from the request URL before redirecting
    url: http://service1:8080 # The URL to which the request will be redirected
  - pattern: /service2/
    strip_prefix: false
    url: http://service2:8080
```

If `strip_prefix` is set to `true`, the `pattern` will be removed from the request URL before the request is redirected. For instance, a request to `http://localhost:8080/service1/...` will be redirected to `http://service1:8080/...`. If `strip_prefix` is `false`, the `pattern` will remain in the redirected URL.

# Deployment Using Docker Compose

You can deploy the API gateway using the following command:

```bash
docker compose up
```

The upstream services must join the external `apigateway-network` that is created with docker-compose.yml. Here is an example docker-compose.yml.

```yaml
version: '3'

services:
  service1:
    container_name: service1
    networks:
      - apigateway-network
    ...
  service2:
    container_name: service2
    networks:
      - apigateway-network
    ...

networks:
  apigateway-network:
    name: apigateway-network
    external: true
```

# Contributing to the Project

Here's how you can contribute:

1. Fork the repository to your own GitHub account.
2. Clone the forked repository to your local machine.
3. Create a new branch for your changes.
4. Make your changes and commit them to your branch.
5. Push your changes to your forked repository.
6. Submit a pull request from your forked repository to our main repository.

Thank you for your interest in contributing!

# Project Licensing

This project is licensed under the MIT License. This means you are free to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the software, subject to the following conditions:

1. The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
2. THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.

For more details, please refer to the [LICENSE](LICENSE) file in the project repository.
