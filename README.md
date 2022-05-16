# envoy-prod-sand-routes

To run the setup using docker-compose, execute the following command from the root of the project repository.

```sh
docker-compose up
```

To get a production token, invoke `/token/prod`.

```sh
curl -X GET \
  'localhost:9095/token/prod'
```

To get a sandbox token, invoke `/token/sand`

```sh
curl -X GET \
  'localhost:9095/token/sand'
```

Invoke the test API endpoint `/api/v1` with production and sandbox tokens to switch between production and sandbox clusters and routes.

Production route match: `/api/v1` ----> `/foo`
Sandbox route match: `/api/v1` ----> `/bar`

## Production endpoint

```sh
curl -X GET \
  'localhost:9095/api/v1' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXBsb3ltZW50IjoiY2x1c3RlclByb2QiLCJleHAiOjE2NTI3Mjg5NDl9.N8nD_4XlhRj-5rUdAKi_gcolscPDyx4MCFIKA39OMf0'
```
Response
```sh
Hello from Foo!!
```

## Sandbox endpoint

```sh
curl -X GET \
  'localhost:9095/api/v1' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXBsb3ltZW50IjoiY2x1c3RlclNhbmQiLCJleHAiOjE2NTI3Mjg5NDl9.44DX4FhcCzYB0_Yi34yyP_ighGtfZogZvEsF9isCmZU'
```
Response
```sh
Hello from Bar!!
```

## envoy route config

```yaml
routes:
    - match:
        prefix: "/api/v1"
        headers:
        - name: x-wso2-cluster
            string_match:
                exact: "clusterSand"
    route:
        prefix_rewrite: /bar
        cluster_header: x-wso2-cluster
    - match:
        prefix: "/api/v1"
    route:
        prefix_rewrite: /foo
        cluster_header: x-wso2-cluster
```
