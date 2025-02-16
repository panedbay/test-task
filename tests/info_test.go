package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func RunInfo(token string) (int, error) {
	client := &http.Client{}
	path := "http://localhost:8080/api/info"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return -1, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, nil
	}

	return resp.StatusCode, nil
}

// Test Info with bought merch and send and received coins
func TestInfo(t *testing.T) {
	payloadFirst := map[string]string{
		"username": "info1",
		"password": "info1",
	}

	jsonData, err := json.Marshal(payloadFirst)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	tokenFirst, code, err := RunAuth(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}

	payloadSecond := map[string]string{
		"username": "info2",
		"password": "info2",
	}

	jsonData, err = json.Marshal(payloadSecond)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	tokenSecond, code, err := RunAuth(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}

	// info1 buys t-shirt
	_, code, err = RunBuy("t-shirt", tokenFirst.Token)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}
	// info1 buys t-shirt

	// info2 sends 10 to info1
	payloadSendFirst := map[string]interface{}{
		"toUser": "info1",
		"amount": 10,
	}

	jsonData, err = json.Marshal(payloadSendFirst)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	_, code, err = RunSendCoin(jsonData, tokenSecond.Token)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}
	// info2 sends 10 to info1

	// info1 sends 10 to info2
	payloadSendSecond := map[string]interface{}{
		"toUser": "info2",
		"amount": 10,
	}

	jsonData, err = json.Marshal(payloadSendSecond)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	_, code, err = RunSendCoin(jsonData, tokenFirst.Token)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}
	// info1 sends 10 to info2

	// get info of info1
	code, err = RunInfo(tokenFirst.Token)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}
	// get info of info1

}
