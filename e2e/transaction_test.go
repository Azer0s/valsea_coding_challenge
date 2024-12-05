package e2e

import (
	"github.com/shopspring/decimal"
	"net/http"
	"slices"
	"testing"
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/enum"
	"valsea_coding_challenge/domain/transactional"
)

//TESTS
// -[x] Create a transaction
// -[x] Get all transactions
// -[x] Create a transaction with an invalid user
// -[x] Create a transaction with an invalid transaction type
// -[x] Create a transaction with an invalid amount

func TestCreateTransactionsAndGetAllTransactions(t *testing.T) {
	user := createUser(transactional.CreateUserRequest{
		Owner:          "Eve",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	t1 := createTransaction(user, transactional.CreateTransactionRequest{
		Amount: decimal.NewFromInt(50),
		Type:   enum.TransactionTypeDeposit,
	}, t)

	t2 := createTransaction(user, transactional.CreateTransactionRequest{
		Amount: decimal.NewFromInt(-25),
		Type:   enum.TransactionTypeWithdrawal,
	}, t)

	transactions := getAllTransactions(user, t)

	if !slices.ContainsFunc(transactions, func(i dto.TransactionDTO) bool {
		return i.Id == t1.Id
	}) {
		t.Errorf("Transaction %s not found", t1.Id)
	}

	if !slices.ContainsFunc(transactions, func(i dto.TransactionDTO) bool {
		return i.Id == t2.Id
	}) {
		t.Errorf("Transaction %s not found", t2.Id)
	}
}

func TestCreateTransactionWithNonexistentUser(t *testing.T) {
	httpReq, err := createRequest("POST", addr+"/accounts/foo/transactions", transactional.CreateTransactionRequest{
		Amount: decimal.NewFromInt(100),
		Type:   enum.TransactionTypeDeposit,
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status code 404, got %d", res.StatusCode)
	}
}

func TestCreateTransactionWithInvalidTransactionType(t *testing.T) {
	u := createUser(transactional.CreateUserRequest{
		Owner:          "Gabe",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	httpReq, err := createRequest("POST", addr+"/accounts/"+u+"/transactions", transactional.CreateTransactionRequest{
		Amount: decimal.NewFromInt(100),
		Type:   enum.TransactionTypeTransferOut,
	})

	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code 400, got %d", res.StatusCode)
	}
}

func expectFailedTransaction(t *testing.T, userId string, req transactional.CreateTransactionRequest) {
	httpReq, err := createRequest("POST", addr+"/accounts/"+userId+"/transactions", req)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code 500, got %d", res.StatusCode)
	}
}

func TestCreateTransactionWithInvalidAmount(t *testing.T) {
	u := createUser(transactional.CreateUserRequest{
		Owner:          "Fiona",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	expectFailedTransaction(t, u, transactional.CreateTransactionRequest{
		Amount: decimal.NewFromInt(0),
		Type:   enum.TransactionTypeDeposit,
	})

	expectFailedTransaction(t, u, transactional.CreateTransactionRequest{
		Amount: decimal.NewFromInt(-1),
		Type:   enum.TransactionTypeDeposit,
	})

	expectFailedTransaction(t, u, transactional.CreateTransactionRequest{
		Amount: decimal.NewFromInt(10),
		Type:   enum.TransactionTypeWithdrawal,
	})

	expectFailedTransaction(t, u, transactional.CreateTransactionRequest{
		Amount: decimal.NewFromInt(-1000),
		Type:   enum.TransactionTypeWithdrawal,
	})
}
