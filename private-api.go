package poloniex

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// GetBalances returns all of available balances.
func (c *Client) GetBalances(ctx context.Context) (*Balances, error) {
	resp, err := c.doPrivateAPIRequest(ctx, "returnBalances", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var ret map[string]string
	if err := c.decodeResponse(resp, &ret, nil); err != nil {
		return nil, err
	}
	if msg, got := ret["error"]; got {
		return nil, errors.New(msg)
	}

	var balance Balances
	balance.Pair = make(map[string]float64)
	for k, v := range ret {
		balance.Pair[k], _ = strconv.ParseFloat(v, 64)
	}

	return &balance, nil
}

// Balances is a pair of symbol with ammount.
type Balances struct {
	Pair map[string]float64
}

// GetCompleteBalances returns all of balances, including available balance,
// balance on orders, and the estimated BTC value of balance.
func (c *Client) GetCompleteBalances(ctx context.Context) (CompleteBalances, error) {
	resp, err := c.doPrivateAPIRequest(ctx, "returnCompleteBalances", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var ret completeBalancesResponse
	if err := c.decodeResponse(resp, &ret, nil); err != nil {
		return nil, err
	}

	if ret.Error != "" {
		return nil, errors.New(ret.Error)
	}
	return ret.CompleteBalances, nil
}

type completeBalancesResponse struct {
	CompleteBalances
	Error string
}

func (c *completeBalancesResponse) UnmarshalJSON(data []byte) error {
	type errorT struct {
		Error string `json:"error,omitempty"`
	}
	var errResp errorT
	if err := json.Unmarshal(data, &errResp); err != nil {
		return err
	}

	if errResp.Error != "" {
		c.Error = errResp.Error
		return nil
	}

	return json.Unmarshal(data, &(c.CompleteBalances))
}

// CompleteBalances is a pair of symbol with Balance.
type CompleteBalances map[string]Balance

// A Balance is including available balance, balance on order, and the
// estimated BTC value of balance.
type Balance struct {
	Available string
	OnOrders  string
	BtcValue  string
}

func (c *Client) Withdraw(ctx context.Context, req WithdrawRequest) error {
	v := url.Values{}
	v.Add("currency", req.Currency)
	v.Add("amount", strconv.FormatFloat(req.Amount, 'f', -1, 64))
	v.Add("address", req.Address)

	resp, err := c.doPrivateAPIRequest(ctx, "withdraw", v)
	if err != nil {
		return err
	}

	var ret withDrawResponse
	if err := c.decodeResponse(resp, &ret, nil); err != nil {
		return err
	}

	if ret.Error != "" {
		return errors.New(ret.Error)
	}

	return nil
}

type WithdrawRequest struct {
	Currency string
	Amount   float64
	Address  string
}

type withDrawResponse struct {
	Error    string
	Response string
}

func (c *Client) Order(ctx context.Context, order OrderRequest) (*Order, error) {
	v := url.Values{}
	v.Add("currencyPair", order.CurrencyPair)
	v.Add("rate", strconv.FormatFloat(order.Rate, 'f', -1, 64))
	v.Add("amount", strconv.FormatFloat(order.Amount, 'f', -1, 64))

	if order.FillOrKill {
		v.Add("fillOrKill", "1")
	}

	if order.ImmediateOrCancel {
		v.Add("immediateOrCancel", "1")
	}

	if order.PostOnly {
		v.Add("postOnly", "1")
	}

	resp, err := c.doPrivateAPIRequest(ctx, string(order.Type), v)
	if err != nil {
		return nil, err
	}

	var ret orderResponse
	if err := c.decodeResponse(resp, &ret, nil); err != nil {
		return nil, err
	}

	if ret.Error != "" {
		return nil, errors.New(ret.Error)
	}

	return &ret.Order, nil
}

type OrderType string

const ORDER_TYPE_BUY OrderType = "buy"
const ORDER_TYPE_SELL OrderType = "sell"

type OrderRequest struct {
	Type              OrderType
	CurrencyPair      string
	Rate              float64
	Amount            float64
	FillOrKill        bool
	ImmediateOrCancel bool
	PostOnly          bool
}

type orderResponse struct {
	Error string `json:"error"`
	Order
}

type Order struct {
	OrderNumber     int            `json:"orderNumber"`
	ResultingTrades []OrderRequest `json:"resultingTrades"`
}

type OrderResult struct {
	Amount  string `json:"amount"`
	Date    string `json:"date"`
	Rate    string `json:"rate"`
	Total   string `json:"total"`
	TradeID string `json:"tradeID"`
	Type    string `json:"type"`
}
