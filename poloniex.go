package poloniex

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"io"
	"os"

	"io/ioutil"
	"sync"
)

func New(apiKey, apiSecret, endpoint string, logger *log.Logger) (*Client, error) {
	if endpoint == "" {
		endpoint = "https://poloniex.com/"
	}
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	return &Client{
		Key:        apiKey,
		Secret:     apiSecret,
		HTTPClient: http.DefaultClient,
		Endpoint:   endpointURL,
		Logger:     logger,
	}, nil
}

type Client struct {
	HTTPClient *http.Client
	Endpoint   *url.URL
	Key        string
	Secret     string
	Logger     *log.Logger
	muNonce    sync.Mutex
}

func (c *Client) decodeResponse(resp *http.Response, ret interface{}, f *os.File) error {
	defer resp.Body.Close()

	if f != nil {
		resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, f))
		defer f.Close()
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(ret); err != nil {
		return err
	}
	return nil
}

func (c *Client) doPrivateAPIRequest(ctx context.Context, command string, values url.Values) (*http.Response, error) {
	c.muNonce.Lock()
	defer c.muNonce.Unlock()

	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	req, err := c.newPrivateAPIRequest(command, nonce, values)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	return c.HTTPClient.Do(req)
}

func (c *Client) newPrivateAPIRequest(command string, nonce string, values url.Values) (*http.Request, error) {
	if values == nil {
		values = url.Values{}
	}

	values.Set("command", command)
	values.Set("nonce", nonce)
	body := values.Encode()

	u := *c.Endpoint
	u.Path = path.Join(c.Endpoint.Path, "tradingApi")

	c.Logger.Println("POST", u.String())
	req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Sign request
	h := hmac.New(sha512.New, []byte(c.Secret))
	h.Write([]byte(body))
	sign := hex.EncodeToString(h.Sum(nil))
	req.Header.Add("Key", c.Key)
	req.Header.Add("Sign", sign)

	return req, nil
}

func (c *Client) newPublicAPIRequest(ctx context.Context, command string, values url.Values) (*http.Request, error) {
	if values == nil {
		values = url.Values{}
	}
	values.Set("command", command)

	u := *c.Endpoint
	u.Path = path.Join(c.Endpoint.Path, "public")
	u.RawQuery = values.Encode()

	c.Logger.Println("GET", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Add("Accept", "application/json")

	return req, nil
}
