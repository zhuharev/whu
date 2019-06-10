package client

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bloom42/rz-go/v2"
	"github.com/bloom42/rz-go/v2/log"
	"github.com/pkg/errors"
	"github.com/zhuharev/whu/domain/update"
)

type Client struct {
	offset     int
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Run(fn func([]byte) error) {
	for t := time.NewTicker(30 * time.Second); ; <-t.C {
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
	err = json.NewDecoder(resp.Body).Decode(&updates)
	if err != nil {
		return errors.Wrap(err, "decode json")
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
