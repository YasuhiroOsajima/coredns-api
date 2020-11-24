# REST API for standalone coredns

## Install

### build

```bash
cd coredns-api/

bash scripts/code_build.sh 
or
bash script/docker_build.sh
```

### set depending OS environment variables
- SERVER  
  This API server's IP address.
- PORT  
  This API server's HTTP listen port.
- CONF_PATH  
  File path of coredns conf.
- HOSTS_DIR  
  Directory path of coredns hosts files.

```bash
vim docker-compose.yml
```

### start

```bash
docker-compose up
```

## Usage

### API doc

Swagger is available on `http://${SERVER}:${PORT}/swagger/index.html`.

#### Add domain

request

```bash
curl -X POST http://127.0.0.1:8080/v1/domains \
-H "Accept: application/json" \
-d '{"domain": "hogehoge.hoge",
     "tenants": ["df397e50-8006-450e-b18b-5c5bd940baff", "02c03bd4-fe2e-45f2-85b6-b535af15215d"]}'
```

response

```text
HTTP/1.1 201 Created
Content-Type: application/json

{"domain": "hogehoge.hoge", "uuid": "aea6cf49-2912-42af-b903-dae1312f64d9", "hosts": [], "tenants": ["df397e50-8006-450e-b18b-5c5bd940baff", "02c03bd4-fe2e-45f2-85b6-b535af15215d"]}
```

#### Update domain

request

```bash
curl -X PATCH http://127.0.0.1:8080/v1/domains/{DOMAIN_UUID} \
-H "Accept: application/json" \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff" \
-d '{"tenants": ["df397e50-8006-450e-b18b-5c5bd940baff", "02c03bd4-fe2e-45f2-85b6-b535af15215d"]}'
```

or

```bash
curl -X PATCH http://127.0.0.1:8080/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID} \
-H "Accept: application/json" \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff" \
-d '{"address": "172.21.1.2"}'
```


#### Delete domain

request

```bash
curl -X DELETE http://127.0.0.1:8080/v1/domains/{DMAIN_UUID} \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff"
```

response

```text
HTTP/1.1 204 No Content
Content-Length: 0
```

#### List domains

request

```bash
curl -X GET http://127.0.0.1:8080/v1/domains \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff"
```

response

```text
HTTP/1.1 200 OK
Content-Type: application/json

{
    "domains": [
        {"domain": "hogehoge.hoge", "uuid": "aea6cf49-2912-42af-b903-dae1312f64d9"},
        {"domain": "fugafuga.hoge", "uuid": "1cf4caeb-f474-44d1-8eda-b9596cc22f00"}
    ]
}
```

#### Get domain

request

```bash
curl -X GET http://127.0.0.1:8080/v1/domains/{DOMAIN_UUID} \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff"
```

response

```text
HTTP/1.1 200 OK
Content-Type: application/json

{
    "domain": "hogehoge.hoge",
    "uuid": "aea6cf49-2912-42af-b903-dae1312f64d9",
    "hosts": [
        {
            "hostname": "hogeserver2",
            "address": "172.21.1.2",
            "uuid": "a51d334d-567c-4566-b1ff-186446403d3a"
        }
    ]
}
```


#### Add host

request

```bash
curl -X POST http://127.0.0.1:8080/v1/domains/{DOMAIN_UUID}/hosts \
-H "Accept: application/json" \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff" \
-d '{"hostname": "hogeserver1", "address": "172.21.1.1"}'
```

response

```text
HTTP/1.1 201 Created
Content-Type: application/json

{
    "domain": "hogehoge.hoge",
    "hosts": {
        "hostname": "hogeserver1",
        "address": "172.21.1.1",
        "uuid": "a51d334d-567c-4566-b1ff-186446403d3a"
    }
}
```

#### Update host

request

```bash
curl -X PATCH http://127.0.0.1:8080/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID} \
-H "Accept: application/json" \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff" \
-d '{"hostname": "hogeserver2"}'
```

or

```bash
curl -X PATCH http://127.0.0.1:8080/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID} \
-H "Accept: application/json" \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff" \
-d '{"address": "172.21.1.2"}'
```

or

```bash
curl -X PATCH http://127.0.0.1:8080/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID} \
-H "Accept: application/json" \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff" \
-d '{"hostname": "hogeserver2", "address": "172.21.1.2"}'
```

response

```text
HTTP/1.1 200 OK
Content-Type: application/json

{
    "domain": "hogehoge.hoge",
    "hosts": {
        "hostname": "hogeserver2",
        "address": "172.21.1.2",
        "uuid": "a51d334d-567c-4566-b1ff-186446403d3a"
    }
}
```

#### Delete host

request

```bash
curl -X DELETE http://127.0.0.1:8080/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID} \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff"
```

response

```text
HTTP/1.1 204 No Content
Content-Length: 0
```

#### Get host

request

```bash
curl -X GET http://127.0.0.1:8080/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID} \
-H "Tenant: df397e50-8006-450e-b18b-5c5bd940baff"
```

response

```text
HTTP/1.1 200 OK
Content-Type: application/json

{
    "domain": "hogehoge.hoge",
    "uuid": "aea6cf49-2912-42af-b903-dae1312f64d9",
    "hosts": [
        {
            "hostname": "hogeserver2",
            "address": "172.21.1.2",
            "uuid": "a51d334d-567c-4566-b1ff-186446403d3a"
        }
    ]
}
```

### DNS query

```bash
dig @127.0.0.1 hogeserver1.hogehoge.hoge
```

### Tenant list

build command

```bash
bash scripts/code_build.sh
```

get tenant list, and its accessible domains.

```bash
bash scripts/tenant_list.sh
```

## Todo
- Use goroutine to update filesystem sequentially.