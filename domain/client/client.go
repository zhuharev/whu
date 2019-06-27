package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/bloom42/rz-go/v2"
	"github.com/bloom42/rz-go/v2/log"
	"github.com/pkg/errors"
	zhuerrors "github.com/zhuharev/errors"
	"github.com/zhuharev/whu/domain/update"
)

type Client struct {
	offset     int
	baseURL    string
	httpClient *http.Client
	interval   int
}

type ClientOpt func(*Client)

func FetchInterval(seconds int) ClientOpt {
	return func(c *Client) {
		c.interval = seconds
	}
}

func New(baseURL string, opts ...ClientOpt) *Client {
	c := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		interval:   15,
	}
	for _, fn := range opts {
		fn(c)
	}
	return c
}

func (c *Client) Run(fn func([]byte) error) {
	for t := time.NewTicker(time.Duration(c.interval) * time.Second); ; <-t.C {
		err := c.doRequest(fn)
		if err != nil {
			log.Error("err do", rz.Err(err))
		}
	}
}

func (c *Client) doRequest(fn func([]byte) error) error {
	resp, err := c.httpClient.Get(c.baseURL + "/updates?offset=" + strconv.Itoa(c.offset))
	if err != nil {
		return errors.Wrap(err, "get http")
	}
	defer resp.Body.Close()
	var updates []update.Update
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ReadResponseBody.New("updates response").String("body", string(data))
	}
	err = json.Unmarshal(data, &updates)
	if err != nil {
		return UnmarshalJSON.New("unmarshal json of body bytes").
			String("body", string(data))
	}
	for _, upd := range updates {
		err := fn(upd.Payload)
		if err != nil {
			return errors.Wrap(err, "callback")
		}
		c.offset = upd.ID
	}
	return nil
}

type ErrorType = zhuerrors.ErrorType

const (
	ReadResponseBody ErrorType = ErrorType(iota)
	UnmarshalJSON
)
