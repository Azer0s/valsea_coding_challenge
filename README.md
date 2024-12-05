# valsea_coding_challenge

A simple backend for a banking system.

## Requirements
- [x] Concurrently support multiple accounts
- [x] Support for deposit and withdrawal
- [x] Support for transfer between accounts
- [x] Testing

## Non-goals (left out because of time constraints)
- [ ] Authentication
- [ ] Authorization
- [ ] Observability

## Usage

```bash
go get ./...
PORT=8080 go run main.go
```

The server will be running on `localhost:8080`.
The endpoints are:
- `GET /accounts`: Get all accounts
- `POST /accounts`: Create a new account
- `GET /accounts/{id}`: Get an account by ID
- `POST /accounts/{id}/transactions`: Create a new transaction for an account
- `GET /accounts/{id}/transactions`: Get all transactions for an account
- `POST /transfer`: Transfer money from one account to another

## Testing

I have implemented all the basic e2e tests for the endpoints (mostly testing busines logic and validation). You can run the tests with the following command:

```bash
go test ./...
```