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

func RunSendCoin(jsonData []byte, token string) (model.SendCoinResponse, int, error) {
	client := &http.Client{}
	path := "http://localhost:8080/api/sendCoin"
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(jsonData))
	if err != nil {
		return model.SendCoinResponse{}, -1, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return model.SendCoinResponse{}, -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.SendCoinResponse{}, -1, err
	}

	if resp.StatusCode != http.StatusOK {
		return model.SendCoinResponse{}, resp.StatusCode, nil
	}

	var result map[string]string
	if e := json.Unmarshal(body, &result); e != nil {
		return model.SendCoinResponse{}, -1, e
	}

	description, ok := result["description"]
	if !ok || description == "" {
		return model.SendCoinResponse{}, -1, errors.New("Failed parsing output of sendCoin_test")
	}

	return model.SendCoinResponse{Desc: description}, resp.StatusCode, nil
}

// Basic send test
func TestSendCoin(t *testing.T) {
	payloadsAuth := make([]map[string]string, 0)

	payloadAuthFirst := map[string]string{
		"username": "sc1",
		"password": "sc1",
	}
	payloadAuthSecond := map[string]string{
		"username": "sc2",
		"password": "sc2",
	}
	payloadsAuth = append(payloadsAuth, payloadAuthFirst, payloadAuthSecond)

	payloadSend := map[string]interface{}{
		"toUser": "sc2",
		"amount": 10,
	}

	tokens := make([]model.AuthResponse, 0)

	for _, payload := range payloadsAuth {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Error marshaling payload: %v", err)
		}

		tok, code, err := RunAuth(jsonData)
		if err != nil {
			t.Fatal(err)
		}
		if code != 200 {
			t.Fatalf("Bad status code: %d", code)
		}
		tokens = append(tokens, tok)
	}

	jsonData, err := json.Marshal(payloadSend)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	_, code, err := RunSendCoin(jsonData, tokens[0].Token)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}
}
