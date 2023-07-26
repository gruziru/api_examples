package main

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

const apiPath = "oboz2-order-client-api/v1/messages"

type Client struct {
	apiKey  string
	Server  string
	APIPath string
	client  *resty.Client
}

func NewClient(server string, key string) *Client {
	return &Client{
		apiKey:  key,
		Server:  server,
		APIPath: apiPath,
		client:  resty.New(),
	}
}

func (c Client) PostTrackingOrders(ords []TrackingOrder) error {
	resp, err := c.client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", c.apiKey).
		SetBody(ords).
		Post(fmt.Sprintf("%s/%s", c.Server, c.APIPath))
	if err != nil {
		return fmt.Errorf("post tracking orders: %w", err)
	}

	log.Printf("HTTP Status: %s | Response: %s", resp.Status(), resp.String())

	return nil
}
