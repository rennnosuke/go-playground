package main

import "net/http"

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
