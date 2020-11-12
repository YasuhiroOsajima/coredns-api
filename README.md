# REST API for standalone coredns

## Add domain

request

```bash
curl -X POST http://127.0.0.1/v1/domains \
-H "Accept: application/json" \
-d '{"domain": "hogehoge.hoge"}'
```

response

```text
HTTP/1.1 201 Created
Content-Type: application/json

{"domain": "hogehoge.hoge", "uuid": "aea6cf49-2912-42af-b903-dae1312f64d9"}
```

## Delete domain

request

```bash
curl -X DELETE http://127.0.0.1/v1/domains/{DMAIN_UUID}
```

response

```text
HTTP/1.1 204 No Content
Content-Length: 0
```

## Get domain

request

```bash
curl -X GET http://127.0.0.1/v1/domains/{DOMAIN_UUID}
```

response

```text
HTTP/1.1 200 OK
Content-Type: application/json

{
    "domains": [
        {"domain": "hogehoge.hoge", "uuid": "aea6cf49-2912-42af-b903-dae1312f64d9"}
    ]
}
```

## List domains

request

```bash
curl -X GET http://127.0.0.1/v1/domains
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

## Add host

request

```bash
curl -X POST http://127.0.0.1/v1/domains/{DOMAIN_UUID}/hosts \
-H "Accept: application/json" \
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

## Update host

request

```bash
curl -X PATCH http://127.0.0.1/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID} \
-H "Accept: application/json" \
-d '{"hostname": "hogeserver2"}'
```

or

```bash
curl -X PATCH http://127.0.0.1/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID} \
-H "Accept: application/json" \
-d '{"address": "172.21.1.2"}'
```

or

```bash
curl -X PATCH http://127.0.0.1/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID} \
-H "Accept: application/json" \
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

## Delete host

request

```bash
curl -X DELETE http://127.0.0.1/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID}
```

response

```text
HTTP/1.1 204 No Content
Content-Length: 0
```

## List hosts

request

```bash
curl -X GET http://127.0.0.1/v1/domains/{DOMAIN_UUID}/hosts/{HOST_UUID}
```

response

```text
HTTP/1.1 200 OK
Content-Type: application/json

{
    "domain": "hogehoge.hoge",
    "hosts": [
        {
            "hostname": "hogeserver2",
            "address": "172.21.1.2",
            "uuid": "a51d334d-567c-4566-b1ff-186446403d3a"
        }
    ]
}
```