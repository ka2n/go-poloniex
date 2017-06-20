package poloniex

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

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

type Balances struct {
	Pair map[string]float64
}

func (c *Client) GetCompleteBalances(ctx context.Context) (CompleteBalances, error) {
	resp, err := c.doPrivateAPIRequest(ctx, "returnCompleteBalances", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var ret map[string]interface{}
	if err := c.decodeResponse(resp, &ret, nil); err != nil {
		return nil, err
	}
	if msg, got := ret["error"]; got {
		return nil, errors.New(msg.(string))
	}

	balances := make(CompleteBalances)
	for k, v := range ret {
		balance := v.(map[string]interface{})
		i, err := strconv.ParseFloat(balance["btcValue"].(string), 64)
		if i == 0 || err != nil {
			continue
		}
		balances[k] = Balance{
			Available: balance["available"].(string),
			OnOrders:  balance["onOrders"].(string),
			BtcValue:  balance["btcValue"].(string),
		}
	}

	return balances, nil
}

type CompleteBalances map[string]Balance

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

	var ret WithDrawResponse
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

type WithDrawResponse struct {
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

	var ret OrderResponse
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

type OrderResponse struct {
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
