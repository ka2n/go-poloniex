package poloniex

import (
	"errors"
	"sync"
	"time"

	"gopkg.in/jcelliott/turnpike.v2"
)

// NewPushAPIClient creates new PushAPIClient
func NewPushAPIClient(endpoint string) *PushAPIClient {
	if endpoint == "" {
		endpoint = "wss://api.poloniex.com:443"
	}

	return &PushAPIClient{
		Endpoint: endpoint,
		subs:     make(map[string]chan Ticker),
	}
}

// PushAPIClient client
type PushAPIClient struct {
	Endpoint string
	client   *turnpike.Client

	tickerSubscribed bool
	subs             map[string]chan Ticker
	mu               sync.Mutex
}

func (c *PushAPIClient) SubscribeTicker(pair string, tc chan Ticker) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.subs[pair] != nil {
		return errors.New("already subscribed")
	}
	c.subs[pair] = tc
	return nil
}

func (c *PushAPIClient) UnsubscribeTicker(pair string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if tc := c.subs[pair]; tc != nil {
		close(tc)
		delete(c.subs, pair)
	}

	if c.tickerSubscribed && len(c.subs) == 0 {
		c.client.Unsubscribe("ticker")
		c.tickerSubscribed = false
	}

	return nil
}

func (c *PushAPIClient) Receive() error {
	if err := c.receive(); err != nil {
		return err
	}
	<-c.client.ReceiveDone
	return nil
}

func (c *PushAPIClient) receive() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.tickerSubscribed {
		return errors.New("already subscribed")
	}

	client, err := turnpike.NewWebsocketClient(turnpike.JSON, c.Endpoint, nil, nil, nil)
	if err != nil {
		return err
	}

	if _, err := client.JoinRealm("realm1", nil); err != nil {
		client.Close()
		return err
	}

	if err := client.Subscribe("ticker", nil, c.receiveTicker); err != nil {
		client.Close()
		return err
	}

	c.client = client
	c.tickerSubscribed = true

	return nil
}

func (c *PushAPIClient) receiveTicker(args []interface{}, kwargs map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	currencyPair := args[0].(string)
	tc, ok := c.subs[currencyPair]
	if !ok {
		return
	}

	ticker := Ticker{}
	ticker.Last,
		ticker.LowestAsk,
		ticker.HighestBid,
		ticker.PercentChange,
		ticker.BaseVolume,
		ticker.QuoteVolume,
		ticker.IsFrozen,
		ticker.High24Hr,
		ticker.Low24Hr =
		args[1].(string),
		args[2].(string),
		args[3].(string),
		args[4].(string),
		args[5].(string),
		args[6].(string),
		int(args[7].(float64)),
		args[8].(string),
		args[9].(string)
	ticker.Time = time.Now()
	tc <- ticker
}

func (c *PushAPIClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.tickerSubscribed {
		return nil
	}

	if c.client != nil {
		c.client.Close()
		c.client = nil
	}

	for k, tc := range c.subs {
		close(tc)
		delete(c.subs, k)
	}

	return nil
}
