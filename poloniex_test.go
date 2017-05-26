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
