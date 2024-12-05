package e2e

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"net/http"
	"slices"
	"testing"
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/transactional"
)

//TESTS
// -[x] Create a user
// -[x] Get a user by id
// -[x] Get all users
// -[x] Find a user that does not exist

func TestCreateUserAndFindById(t *testing.T) {
	id := createUser(transactional.CreateUserRequest{
		Owner:          "Alice",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	u := getUserById(id, t)
	if u.Name != "Alice" {
		t.Fatalf("Expected user name Alice, got %s", u.Name)
	}
}

func TestFindNonexistentUser(t *testing.T) {
	res, err := http.Get(addr + "/accounts/123")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status code 404, got %d", res.StatusCode)
	}
}

func TestGetAllUsers(t *testing.T) {
	createUser(transactional.CreateUserRequest{
		Owner:          "Bob",
		InitialBalance: decimal.NewFromInt(100),
	}, t)
	createUser(transactional.CreateUserRequest{
		Owner:          "Charlie",
		InitialBalance: decimal.NewFromInt(100),
	}, t)
	createUser(transactional.CreateUserRequest{
		Owner:          "David",
		InitialBalance: decimal.NewFromInt(100),
	}, t)

	res, err := http.Get(addr + "/accounts")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	var users []dto.UserDTO
	err = json.NewDecoder(res.Body).Decode(&users)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.ContainsFunc(users, func(u dto.UserDTO) bool {
		return u.Name == "Bob"
	}) {
		t.Fatal("Bob not found")
	}

	if !slices.ContainsFunc(users, func(u dto.UserDTO) bool {
		return u.Name == "Charlie"
	}) {
		t.Fatal("Charlie not found")
	}

	if !slices.ContainsFunc(users, func(u dto.UserDTO) bool {
		return u.Name == "David"
	}) {
		t.Fatal("David not found")
	}
}
