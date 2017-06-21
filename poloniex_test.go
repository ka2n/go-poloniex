package poloniex

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestClient_GetCurrencies(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	client, err := New("", "", "", logger)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	ret, err := client.GetCurrencies(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if ret == nil {
		t.Fatal("ret should be non nil")
	}
}

func TestClient_GetBalances(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	client, err := NewPrivateClient(logger)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	ret, err := client.GetBalances(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if ret == nil {
		t.Fatal("ret should be non nil")
	}
}

func TestClient_GetCompleteBalances(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	client, err := NewPrivateClient(logger)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	ret, err := client.GetCompleteBalances(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if ret == nil {
		t.Fatal("ret should be non nil")
	}
}

func NewPrivateClient(logger *log.Logger) (*Client, error) {
	return New("POLONIEX_KEY", "POLONIEX_SECRET", "", logger)
}
