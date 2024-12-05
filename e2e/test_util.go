package e2e

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"strconv"
	"testing"
	"valsea_coding_challenge/cmd"
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/transactional"
	"valsea_coding_challenge/util"
)

var addr string

func getFreePort() int {
	if a, err := net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		if l, err := net.ListenTCP("tcp", a); err == nil {
			defer func(l *net.TCPListener) {
				err := l.Close()
				if err != nil {
					panic(err)
				}
			}(l)
			return l.Addr().(*net.TCPAddr).Port
		}
	}

	panic("Failed to get a free port")
}

func init() {
	port := getFreePort()
	addr = "http://localhost:" + strconv.Itoa(port)

	startingChan := make(chan struct{})
	go cmd.StartServer(&util.Config{
		Port:     port,
		LogLevel: zapcore.DebugLevel,
	}, startingChan)
	<-startingChan
}

func createRequest(method, url string, body interface{}) (*http.Request, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		panic(err)
	}

	return req, nil
}

func createUser(req transactional.CreateUserRequest, t *testing.T) string {
	httpReq, err := createRequest("POST", addr+"/accounts", req)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 201, got %d", res.StatusCode)
	}

	var user dto.UserDTO
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		t.Fatal(err)
	}

	return user.Id
}

func createTransaction(userId string, req transactional.CreateTransactionRequest, t *testing.T) *dto.TransactionDTO {
	httpReq, err := createRequest("POST", addr+"/accounts/"+userId+"/transactions", req)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	var transaction dto.TransactionDTO
	err = json.NewDecoder(res.Body).Decode(&transaction)
	if err != nil {
		t.Fatal(err)
	}

	return &transaction
}

func getAllTransactions(userId string, t *testing.T) []dto.TransactionDTO {
	httpReq, err := createRequest("GET", addr+"/accounts/"+userId+"/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	var transactions []dto.TransactionDTO
	err = json.NewDecoder(res.Body).Decode(&transactions)
	if err != nil {
		t.Fatal(err)
	}

	return transactions
}

func getUserById(id string, t *testing.T) dto.UserDTO {
	res, err := http.Get(addr + "/accounts/" + id)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	var user dto.UserDTO
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func createTransfer(req transactional.CreateTransferRequest, t *testing.T) []dto.TransactionDTO {
	httpReq, err := createRequest("POST", addr+"/transfer", req)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	var transactions []dto.TransactionDTO
	err = json.NewDecoder(res.Body).Decode(&transactions)
	if err != nil {
		t.Fatal(err)
	}

	return transactions
}
