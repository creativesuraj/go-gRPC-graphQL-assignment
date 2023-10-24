# Aeon Take Home Assignment

Develop a Golang application using gRPC to serve data about books. The assignment should showcase
your skills with gRPC, GraphQL, and Golang best practices.

## Installation

```bash
  go mod download
```

## Set Environment Variables

Set the following environment variables for the grpc-server

```bash
  export DB_URI=
  export DB_NAME=
  export GRPC_PORT=
```

Set the following environment variables for the gql-server

```bash
  export GRPC_SERVER_URL=
  export GQL_SERVER_PORT=
```

## Run grpc-server

Go inside `grpc-server` and run

```
  go run server.go
```

## Run gql-server

Go inside `gql-server` and run

```
  go run server.go
```

## Make GraphQL Requests

The GraphQL playground is running on

```
http://localhost:{GQL_SERVER_PORT}/play
```