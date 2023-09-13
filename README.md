# Smartway test

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/18839470-3aee72e0-6b99-454a-8ddf-a4ed483b2444?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D18839470-3aee72e0-6b99-454a-8ddf-a4ed483b2444%26entityType%3Dcollection%26workspaceId%3Dae41c6ce-f7e9-4fdd-9855-aaa4689acb54)

## Operations

### Add airline

```shell
curl --location 'localhost:8080/v1/airlines/S7' \
--header 'Content-Type: application/json' \
--data '{
    "name": "S7"
}'
```

### Delete an airline by ID

```shell
curl --location --request DELETE 'localhost:8080/v1/airlines/S7'
```

### Add provider

```shell
curl --location 'localhost:8080/v1/providers/RS' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Red Star"
}'
```

### Delete a provider by ID

```shell
curl --location --request DELETE 'localhost:8080/v1/providers/RS'
```

### Change providers (ID list) for airline (code)

```shell
curl --location --request PUT 'localhost:8080/v1/airlines/S7/providers' \
--header 'Content-Type: application/json' \
--data '{
    "providers": ["RS"]
}'
```

### Add search schema

```shell
curl --location 'localhost:8080/v1/schemas/42' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Demo"
}'
```

### Find search schema by name

```shell
curl --location 'localhost:8080/v1/schemas/find?name=Primary'
```

Example response...

```json
{
    "name": "Primary",
    "id": 1,
    "providers": [
        "AA",
        "IF"
    ]
}
```

### Change search schema by ID

```shell
curl --location --request PUT 'localhost:8080/v1/schemas/42' \
--header 'Content-Type: application/json' \
--data '{
    "providers": ["RS"]
}'
```

### Delete search schema

> Will fail if the schema is currently assigned to at least one account

```shell
curl --location --request DELETE 'localhost:8080/v1/schemas/42'
```

### Add account

```shell
curl --location --request POST 'localhost:8080/v1/accounts/44'
```

### Change assigned schema

```shell
curl --location --request PUT 'localhost:8080/v1/accounts/44/schema' \
--header 'Content-Type: application/json' \
--data '{
    "id": 42
}'
```

### Delete account

```shell
curl --location --request DELETE 'localhost:8080/v1/accounts/44'
```

### Get airlines by account ID

```shell
curl --location 'localhost:8080/v1/accounts/1/airlines'
```

Example response...


```json
[
    {
        "code": "SU",
        "name": "Аэрофлот"
    },
    {
        "code": "S7",
        "name": "S7"
    }
]
```

### Get airlines by provider ID

```shell
curl --location 'localhost:8080/v1/providers/RS/airlines'
```
