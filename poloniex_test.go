package poloniex

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var logger = log.New(ioutil.Discard, "", 0)

func TestClient_GetCurrencies(t *testing.T) {
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
	server, client := newMockedClient([]mockResponse{
		{Action: "returnBalances", PrivateAPI: true, BodyFile: "testdata/returnBalances.json"},
	})
	defer server.Close()

	ctx := context.Background()
	ret, err := client.GetBalances(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if ret == nil {
		t.Fatal("ret should be non nil")
	}

	expect := &Balances{map[string]float64{
		"LTC": 1000.11111111,
		"BTC": 0.000000003,
	}}

	if !reflect.DeepEqual(expect, ret) {
		t.Errorf("expect: %+v, got: %+v", expect, ret)
	}
}

func TestClient_GetCompleteBalances(t *testing.T) {
	server, client := newMockedClient([]mockResponse{
		{Action: "returnCompleteBalances", PrivateAPI: true, BodyFile: "testdata/returnCompleteBalances.json"},
	})
	defer server.Close()

	ctx := context.Background()
	ret, err := client.GetCompleteBalances(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if ret == nil {
		t.Fatal("ret should be non nil")
	}

	expect := CompleteBalances(map[string]Balance{
		"LTC": Balance{
			Available: "5.015",
			OnOrders:  "1.0025",
			BtcValue:  "0.078",
		},
		"NXT": Balance{
			Available: "0",
			OnOrders:  "0",
			BtcValue:  "0",
		},
	})

	if !reflect.DeepEqual(expect, ret) {
		t.Errorf("expect: %+v, got: %+v", expect, ret)
	}
}

type mockResponse struct {
	Action     string
	PrivateAPI bool
	BodyFile   string
	Status     int
}

func newMockedClient(mocks []mockResponse) (*httptest.Server, *Client) {
	muxAPI := http.NewServeMux()
	testAPIServer := httptest.NewServer(muxAPI)

	public := "/public"
	private := "/tradingApi"

	for _, mock := range mocks {
		var req string

		if mock.PrivateAPI {
			req = private
		} else {
			req = public
		}

		muxAPI.HandleFunc(req, func(w http.ResponseWriter, r *http.Request) {
			if mock.Status != 0 {
				w.WriteHeader(mock.Status)
			}
			http.ServeFile(w, r, mock.BodyFile)
		})
	}

	client, _ := New("", "", testAPIServer.URL, nil)
	return testAPIServer, client
}
