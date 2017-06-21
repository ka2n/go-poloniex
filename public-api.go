package poloniex

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// GetTickers returns latest tickers
func (c *Client) GetTickers(ctx context.Context) (*Tickers, error) {
	req, err := c.newPublicAPIRequest(ctx, "returnTicker", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var ret Tickers
	if err := c.decodeResponse(resp, &ret.Pair, nil); err != nil {
		return nil, err
	}

	return &ret, nil
}

// Tickers holds map of Ticker
type Tickers struct {
	Pair map[string]Ticker
}

// Ticker represent ticker
type Ticker struct {
	Last          string    `json:"last"`
	LowestAsk     string    `json:"lowestAsk"`
	HighestBid    string    `json:"highestBid"`
	PercentChange string    `json:"percentChange"`
	BaseVolume    string    `json:"baseVolume"`
	QuoteVolume   string    `json:"quoteVolume"`
	IsFrozen      int       `json:"isFrozen,string"`
	High24Hr      string    `json:"high24hr"`
	Low24Hr       string    `json:"low24hr"`
	Time          time.Time `json:"-"`
}

// GetCurrencies returns all currencies available
func (c *Client) GetCurrencies(ctx context.Context) (*Currencies, error) {
	req, err := c.newPublicAPIRequest(ctx, "returnCurrencies", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var ret Currencies
	if err := c.decodeResponse(resp, &ret.Pair, nil); err != nil {
		return nil, err
	}

	return &ret, nil
}

// Currency represent currency profile
type Currency struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	DepositAddress string  `json:"depositAddress"`
	TxFee          float64 `json:"txFee,string"`
	MinConf        int     `json:"minConf"`
	Disabled       int     `json:"disabled"`
	Frozen         int     `json:"frozen"`
	Delisted       int     `json:"delisted"`
}

// Currencies holds map of Currency
type Currencies struct {
	Pair map[string]Currency
}
