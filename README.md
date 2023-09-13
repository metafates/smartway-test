# Smartway test

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/18839470-3aee72e0-6b99-454a-8ddf-a4ed483b2444?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D18839470-3aee72e0-6b99-454a-8ddf-a4ed483b2444%26entityType%3Dcollection%26workspaceId%3Dae41c6ce-f7e9-4fdd-9855-aaa4689acb54)

<!--toc:start-->
- [Smartway test](#smartway-test)
  - [Running](#running)
  - [Operations](#operations)
    - [Add airline](#add-airline)
    - [Delete an airline by ID](#delete-an-airline-by-id)
    - [Add provider](#add-provider)
    - [Delete a provider by ID](#delete-a-provider-by-id)
    - [Change providers (ID list) for airline (code)](#change-providers-id-list-for-airline-code)
    - [Add search schema](#add-search-schema)
    - [Find search schema by name](#find-search-schema-by-name)
    - [Change search schema by ID](#change-search-schema-by-id)
    - [Delete search schema](#delete-search-schema)
    - [Add account](#add-account)
    - [Change assigned schema](#change-assigned-schema)
    - [Delete account](#delete-account)
    - [Get airlines by account ID](#get-airlines-by-account-id)
    - [Get airlines by provider ID](#get-airlines-by-provider-id)
<!--toc:end-->

Project structure is based on [evrone/go-clean-template](https://github.com/evrone/go-clean-template)

## Running

This project uses [Mage](https://magefile.org/) as the build tool

<details>
<summary>Why?</summary>

From the [Mage](https://magefile.org/) website...

> Makefiles are hard to read and hard to write. Mostly because makefiles are
> essentially fancy bash scripts with significant white space and
> additional make-related syntax.
>
> Mage lets you have multiple magefiles, name your magefiles whatever
> you want, and they’re easy to customize for multiple operating systems.
> Mage has no dependencies (aside from go) and runs just fine on all major
> operating systems, whereas make generally uses bash which is not well
> supported on Windows. Go is superior to bash for any non-trivial task
> involving branching, looping, anything that’s not just straight line
> execution of commands. And if your project is written in Go, why
> introduce another language as idiosyncratic as bash?
> Why not use the language your contributors are already comfortable with?

</details>

To start the server run

```shell
mage full
```

This will spin up...

- Server inside a docker container on **port 8080**
- Spin up auxiliary containers with
    - PostgreSQL - The database
    - [PGWeb](https://github.com/sosedoff/pgweb) - Web UI for Postgres
- Run [migrations](./migrations) using [pressly/goose](https://github.com/pressly/goose)
- Load [example data](./example-data.sql) into the DB

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
