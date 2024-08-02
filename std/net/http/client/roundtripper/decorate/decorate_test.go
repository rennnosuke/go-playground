package decorate

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var cli = http.Client{
	Transport: &DecorateRoundTripper{
		base: http.DefaultTransport,
		before: func(r *http.Request) error {
			slog.InfoContext(r.Context(), "before request")
			return nil
		},
		after: func(r *http.Response) error {
			slog.InfoContext(r.Request.Context(), "after response")
			return nil
		},
	},
	Timeout: 30,
}

func TestDecorateRoundTripper_RoundTrip(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	tsURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		base   http.RoundTripper
		before func(*http.Request) error
		after  func(*http.Response) error
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr error
	}{
		{
			name: "no decoration",
			fields: fields{
				base: http.DefaultTransport,
			},
			args: args{
				req: &http.Request{
					Method:     "GET",
					URL:        tsURL,
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     make(http.Header),
					Host:       tsURL.Host,
				},
			},
			want: &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"Content-Length": {"0"},
					"Date":           []string{time.Now().UTC().Format(http.TimeFormat)},
				},
				Body: http.NoBody,
				Request: &http.Request{
					Method:     "GET",
					URL:        tsURL,
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     http.Header{},
					Host:       tsURL.Host,
				},
			},
			wantErr: nil,
		},
		{
			name: "preprocess",
			fields: fields{
				base: http.DefaultTransport,
				before: func(r *http.Request) error {
					r.Header.Set("X-Test", "test")
					return nil
				},
			},
			args: args{
				req: &http.Request{
					Method:     "GET",
					URL:        tsURL,
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     make(http.Header),
					Host:       tsURL.Host,
				},
			},
			want: &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"Content-Length": {"0"},
					"Date":           []string{time.Now().UTC().Format(http.TimeFormat)},
				},
				Body: http.NoBody,
				Request: &http.Request{
					Method:     "GET",
					URL:        tsURL,
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header: http.Header{
						// add header in preprocess
						"X-Test": {"test"},
					},
					Host: tsURL.Host,
				},
			},
			wantErr: nil,
		},
		{
			name: "postprocess",
			fields: fields{
				base: http.DefaultTransport,
				after: func(r *http.Response) error {
					r.Header.Set("X-Test", "test")
					return nil
				},
			},
			args: args{
				req: &http.Request{
					Method:     "GET",
					URL:        tsURL,
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     make(http.Header),
					Host:       tsURL.Host,
				},
			},
			want: &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"Content-Length": {"0"},
					"Date":           []string{time.Now().UTC().Format(http.TimeFormat)},

					// add header in postprocess
					"X-Test": {"test"},
				},
				Body: http.NoBody,
				Request: &http.Request{
					Method:     "GET",
					URL:        tsURL,
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     http.Header{},
					Host:       tsURL.Host,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DecorateRoundTripper{
				base:   tt.fields.base,
				before: tt.fields.before,
				after:  tt.fields.after,
			}
			got, err := d.RoundTrip(tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := cmp.Options{
				cmpopts.IgnoreUnexported(http.Request{}),
			}
			if diff := cmp.Diff(tt.want, got, opts); diff != "" {
				t.Errorf("RoundTrip() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
