package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/panedbay/test-task/model"
)

func RunAuth(jsonData []byte) (model.AuthResponse, int, error) {
	resp, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.AuthResponse{}, -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.AuthResponse{}, -1, err
	}

	if resp.StatusCode != http.StatusOK {
		return model.AuthResponse{}, resp.StatusCode, nil
	}

	var result map[string]string
	if e := json.Unmarshal(body, &result); e != nil {
		return model.AuthResponse{}, -1, e
	}

	token, ok := result["token"]
	if !ok || token == "" {
		return model.AuthResponse{}, -1, errors.New("Failed parsing output of auth_test")
	}

	return model.AuthResponse{Token: token}, resp.StatusCode, nil
}

// Basic auth test
func TestAuthFirst(t *testing.T) {
	payload := map[string]string{
		"username": "auth1",
		"password": "auth1",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	_, code, err := RunAuth(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}
}

// Second auth with same data
func TestAuthSecond(t *testing.T) {
	payload := map[string]string{
		"username": "auth2",
		"password": "auth2",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	_, code, err := RunAuth(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}

	_, code, err = RunAuth(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}
}

// Auth with incorrect password
func TestAuthIncorrectPass(t *testing.T) {
	payload := map[string]string{
		"username": "auth3",
		"password": "auth3",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	_, code, err := RunAuth(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}

	payload = map[string]string{
		"username": "auth4",
		"password": "auth4",
	}

	jsonData, err = json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	// Check that incorrect password is not accepted
	_, code, err = RunAuth(jsonData)
	if err != nil && code != 401 {
		t.Fatal(err)
	}

}
