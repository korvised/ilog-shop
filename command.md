### Open in bash command line

## Generate gRPC code in auth module
```shell
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./modules/auth/authPb/authPb.proto
```

## Generate gRPC code in player module
```shell
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./modules/player/playerPb/playerPb.proto
```

## Generate gRPC code in item module
```shell
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./modules/item/itemPb/itemPb.proto
```

## Generate gRPC code in inventory module
```shell
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./modules/inventory/inventoryPb/inventoryPb.proto
```