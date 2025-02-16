package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/panedbay/test-task/model"
)

func RunBuy(item, token string) (model.BuyItemResponse, int, error) {
	client := &http.Client{}
	path := fmt.Sprintf("http://localhost:8080/api/buy/%s", item)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return model.BuyItemResponse{}, -1, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return model.BuyItemResponse{}, -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.BuyItemResponse{}, -1, err
	}

	if resp.StatusCode != http.StatusOK {
		return model.BuyItemResponse{}, resp.StatusCode, nil
	}

	var result map[string]string
	if e := json.Unmarshal(body, &result); e != nil {
		return model.BuyItemResponse{}, -1, e
	}

	description, ok := result["description"]
	if !ok || description == "" {
		return model.BuyItemResponse{}, -1, errors.New("Failed parsing output of buy_test")
	}

	return model.BuyItemResponse{Desc: description}, resp.StatusCode, nil
}

// Buy first part of items that are available
func TestBuyItemsP1(t *testing.T) {
	payload := map[string]string{
		"username": "u1",
		"password": "p1",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	token, code, err := RunAuth(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}

	items := []string{
		"t-shirt",
		"cup",
		"book",
		"pen",
		"powerbank",
		"hoody",
	}

	for _, item := range items {
		_, code, err := RunBuy(item, token.Token)
		if err != nil {
			t.Fatal(err)
		}
		if code != 200 {
			t.Fatalf("Bad status code: %d", code)
		}
	}
}

// Buy second part of items that are available
func TestBuyItemsP2(t *testing.T) {
	payload := map[string]string{
		"username": "u2",
		"password": "p2",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	token, code, err := RunAuth(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}

	items := []string{
		"umbrella",
		"socks",
		"wallet",
		"pink-hoody",
	}

	for _, item := range items {
		_, code, err := RunBuy(item, token.Token)
		if err != nil {
			t.Fatal(err)
		}
		if code != 200 {
			t.Fatalf("Bad status code: %d", code)
		}
	}
}

// Buy non-existent item
func TestBuyItemsFaulty(t *testing.T) {
	payload := map[string]string{
		"username": "u3",
		"password": "p3",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error marshaling payload: %v", err)
	}

	token, code, err := RunAuth(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if code != 200 {
		t.Fatalf("Bad status code: %d", code)
	}

	items := []string{
		"fridge",
	}

	for _, item := range items {
		_, code, err := RunBuy(item, token.Token)
		if err != nil {
			t.Fatal(err)
		}
		if code == 200 {
			t.Fatalf("Should not return OK code: %d", code)
		}
	}
}

// Buy item without JWT token
func TestBuyItemsNoToken(t *testing.T) {
	token := model.AuthResponse{Token: "qwe"}

	items := []string{
		"t-shirt",
	}

	for _, item := range items {
		_, code, err := RunBuy(item, token.Token)
		if err != nil {
			t.Fatal(err)
		}
		if code == 200 {
			t.Fatalf("Should not return OK code: %d", code)
		}
	}
}
