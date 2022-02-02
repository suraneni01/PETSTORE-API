# Pet Store Service

Service is a lightweight proxy fot pet-store API example that additonaly handles input validation.

### Project layout
```
cmd                  | app configuration and running

  internal           | app internal logic

    petstore         | stores client to pet-store API logic

      client_mock.go | auto-generated mock for API client
      client_test.go | unit-tests for API client
      client.go      | HTTP API client implementation
      interface.go   | declares client external contract
      models.go      | declates models and types

    server           | stores service logic

      models.go      | declares request/response models
      server_test.go | unit-tests for service
      server.go      | HTTP service implementation


vendor               \ 
go.mod                | dependencies
go.sum               /

```

### Service Methods

```
GET  http://{HOST}/pets     | Return list of pets
GET  http://{HOST}/pets/id  | Return pet by id
POST http://{HOST}/pets     | Adds new pet
```

### Run server
```go
go run cmd/main.go
```

### Run tests
```go
go test -v ./...
```