# apigateway

[![Test](https://github.com/ermanimer/apigateway/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/ermanimer/apigateway/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/ermanimer/apigateway/graph/badge.svg?token=rbFp8CZIRk)](https://codecov.io/gh/ermanimer/apigateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/ermanimer/apigateway)](https://goreportcard.com/report/github.com/ermanimer/apigateway)

apigateway development environment yada sandboxlarda kullanilmak uzere tasarlanmis bir api gatewaydir. Dizayn esnasinda asagidaki kararlari takip eder:

 - Konfigurasyonu basittir.
 - Docker compose icerisinden direkt olarak kullanilabilir.
 - Load balancing yapmaz.
 - Request ve response'lari degistirmez.
 - Standard kutuphanede bulunan reverse proxy'yi kullanir, verimliligi yuksektir.

 # Docker Compose Ile Kullanmak

docker-compose.yml

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

config.yml

```yaml
upstreams:
  - pattern: /service1/
    strip_prefix: true
    url: http://service1:8080
```

Bu konfigurasyona sahip docker-compose ile api gateway'e gelen /service1/ pattern'i ile baslayan tum istekler, service1'e pattern strip edilerek iletilecektir. Ornek:

```
http://localhost:80/service1/health-check -> http://service1/8080/health-check:
```

Eger strip_prefix false olsaydi api gateway'2 e gelen yukaridaki istek pattern strip edilmeden su sekilde service1'e iletilecekti.

```
http://localhost:80/service1/health-check -> http://service1/service1/8080/health-check:
```

# Konfigurasyon

Yukaridaki ornekte goruldugu gibi sadece upstream konfigurasyonu yeterlidir. Asagida server'a ait varsayilan konfigurasyon gosterilmistir. Istenildigi durumda server da konfigure edilebilir.

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

Go ile projeyi build etmek icin proje dizininde asagidaki komutlari calistirabilirsiniz.

**build:**

```bash
go build ./cmd/apigateway
```

**run:**

```bash
./apigateway
```

Not: apigateway bulundugu dizinde config.yml dosyasini arayacaktir.

# Katilim

Bir issue acin ve istediginiz degisiklikler ve ozelliklere ortak karar verelim. Sonrasinda fork'ladiginiz repository'deki main branch'tan ana repository'deki main branch'e pull request acarak ilerleyebiliriz.

# Lisans

Bu repo MIT lisansina sahiptir.