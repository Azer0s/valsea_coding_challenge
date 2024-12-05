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
// -[x] Transfer money
// -[x] Transfer money with insufficient funds
// -[x] Transfer money with invalid user
// -[x] Transfer money to a user that does not exist

func transferAndCheck(u1 string, u2 string, amount decimal.Decimal, t *testing.T) {
	u1UserPrev := getUserById(u1, t)
	u2UserPrev := getUserById(u2, t)

	t1 := createTransfer(transactional.CreateTransferRequest{
		FromUserId: u1,
		ToUserId:   u2,
		Amount:     amount,
	}, t)

	var t1From dto.TransactionDTO
	var t1To dto.TransactionDTO
	for _, tr := range t1 {
		if tr.TransactionType == enum.TransactionTypeTransferIn {
			t1To = tr
		}
		if tr.TransactionType == enum.TransactionTypeTransferOut {
			t1From = tr
		}
	}

	u1Transactions := getAllTransactions(u1, t)
	u2Transactions := getAllTransactions(u2, t)

	if !slices.ContainsFunc(u1Transactions, func(i dto.TransactionDTO) bool {
		return i.Id == t1From.Id
	}) {
		t.Errorf("Transaction %s not found", t1From.Id)
	}

	if !slices.ContainsFunc(u2Transactions, func(i dto.TransactionDTO) bool {
		return i.Id == t1To.Id
	}) {
		t.Errorf("Transaction %s not found", t1To.Id)
	}

	u1User := getUserById(u1, t)
	u2User := getUserById(u2, t)

	if u1User.Balance != u1UserPrev.Balance-amount.InexactFloat64() {
		t.Errorf("Expected balance %f, got %f", u1UserPrev.Balance-amount.InexactFloat64(), u1User.Balance)
	}

	if u2User.Balance != u2UserPrev.Balance+amount.InexactFloat64() {
		t.Errorf("Expected balance %f, got %f", u2UserPrev.Balance+amount.InexactFloat64(), u2User.Balance)
	}
}

func TestTransferAndGetTransactions(t *testing.T) {
	u1 := createUser(transactional.CreateUserRequest{
		Owner:          "User1",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	u2 := createUser(transactional.CreateUserRequest{
		Owner:          "User2",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	transferAndCheck(u1, u2, decimal.NewFromInt(10), t)
	transferAndCheck(u2, u1, decimal.NewFromInt(10), t)
}

func TestTransferWithInsufficientFunds(t *testing.T) {
	u1 := createUser(transactional.CreateUserRequest{
		Owner:          "User10",
		InitialBalance: decimal.NewFromInt(100),
	}, t)
	u2 := createUser(transactional.CreateUserRequest{
		Owner:          "User20",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	httpReq, err := createRequest("POST", addr+"/transfer", transactional.CreateTransferRequest{
		FromUserId: u1,
		ToUserId:   u2,
		Amount:     decimal.NewFromInt(200),
	})

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

	u1User := getUserById(u1, t)
	if u1User.Balance != 100 {
		t.Fatalf("Expected balance 100, got %f", u1User.Balance)
	}

	u2User := getUserById(u2, t)
	if u2User.Balance != 100 {
		t.Fatalf("Expected balance 100, got %f", u2User.Balance)
	}
}

func TestTransferToNonexistentUser(t *testing.T) {
	u1 := createUser(transactional.CreateUserRequest{
		Owner:          "User100",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	httpReq, err := createRequest("POST", addr+"/transfer", transactional.CreateTransferRequest{
		FromUserId: u1,
		ToUserId:   "invalid",
		Amount:     decimal.NewFromInt(10),
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

	u1User := getUserById(u1, t)
	if u1User.Balance != 100 {
		t.Fatalf("Expected balance 100, got %f", u1User.Balance)
	}
}

func TestTransferFromNonexistentUser(t *testing.T) {
	u1 := createUser(transactional.CreateUserRequest{
		Owner:          "User200",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	httpReq, err := createRequest("POST", addr+"/transfer", transactional.CreateTransferRequest{
		FromUserId: "invalid",
		ToUserId:   u1,
		Amount:     decimal.NewFromInt(10),
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

	u1User := getUserById(u1, t)
	if u1User.Balance != 100 {
		t.Fatalf("Expected balance 100, got %f", u1User.Balance)
	}
}
