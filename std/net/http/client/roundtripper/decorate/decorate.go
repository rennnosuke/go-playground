package main

import (
	"bytes"
	"io"
	"net/http"
)

type DecorateRoundTripper struct {
	base   http.RoundTripper
	before func(*http.Request) error
	after  func(*http.Response) error
}

func (d *DecorateRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if d.before != nil {
		if err := d.before(req); err != nil {
			return nil, err
		}
	}
	resp, err := d.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if d.after != nil {
		if err := d.after(resp); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

type ConsumableDecorateRoundTripper struct {
	base   http.RoundTripper
	before func(*http.Request) error
	after  func(*http.Response) error
}

func (c *ConsumableDecorateRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.before != nil {
		rb, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		if err := c.before(req); err != nil {
			return nil, err
		}
		req.Body = rb
	}
	resp, err := c.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if c.after != nil {
		rb, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body = io.NopCloser(bytes.NewReader(rb))
		if err := c.after(resp); err != nil {
			return nil, err
		}
	}
	return resp, nil
}
