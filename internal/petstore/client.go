package petstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var _ Interface = (*Client)(nil)

type Client struct {
	url  string
	http *http.Client
}

func NewClient(url string, timeout time.Duration) *Client {
	return &Client{
		url:  url,
		http: &http.Client{Timeout: timeout},
	}
}

func (c *Client) Add(t PetType, price float32) error {
	requestPayload := map[string]interface{}{
		"type":  t,
		"price": price,
	}

	body, err := json.Marshal(requestPayload)
	if err != nil {
		return fmt.Errorf("marshal request payload failed: %w", err)
	}

	url := c.url + "/pets"
	resp, err := c.http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d: %s", resp.StatusCode, body)
	}
	return nil
}

func (c *Client) Get(id uint64) (*Pet, error) {
	url := c.url + "/pets/" + strconv.FormatUint(id, 10)
	resp, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, body)
	}

	var pet Pet
	err = json.Unmarshal(body, &pet)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body failed: %w", err)
	}

	return &pet, nil
}

func (c *Client) List() ([]*Pet, error) {
	url := c.url + "/pets"
	resp, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, body)
	}

	var pets []*Pet
	err = json.Unmarshal(body, &pets)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response body failed: %w", err)
	}

	return pets, nil
}
