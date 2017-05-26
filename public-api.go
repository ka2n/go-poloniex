package poloniex

import (
	"context"
	"errors"
	"net/http"
)

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

type Tickers struct {
	Pair map[string]Ticker
}

type Ticker struct {
	Last          float64 `json:"last,string"`
	LowestAsk     float64 `json:"lowestAsk,string"`
	HighestBid    float64 `json:"highestBid,string"`
	PercentChange float64 `json:"percentChange,string"`
	BaseVolume    float64 `json:"baseVolume,string"`
	QuoteVolume   float64 `json:"quoteVolume,string"`
	IsFrozen      int     `json:"isFrozen,string"`
	High24Hr      float64 `json:"high24hr,string"`
	Low24Hr       float64 `json:"low24hr,string"`
}

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

type Currencies struct {
	Pair map[string]Currency
}
