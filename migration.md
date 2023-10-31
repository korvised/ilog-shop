# Migration Database

## Auth Database

```shell
go run ./pkg/database/script/migration.go ./env/dev/.env.auth
```

## Player Database

```shell
go run ./pkg/database/script/migration.go ./env/dev/.env.player
```

## Item Database

```shell
go run ./pkg/database/script/migration.go ./env/dev/.env.item
```

## Inventory Database

```shell
go run ./pkg/database/script/migration.go ./env/dev/.env.inventory
```

## Payment Database

```shell
go run ./pkg/database/script/migration.go ./env/dev/.env.payment
```