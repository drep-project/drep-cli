// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package component

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/drep-project/drepcli/log"
	rpcTypes "github.com/drep-project/drepcli/rpc/types"
	"github.com/rs/cors"
)

var nullAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:0")

type httpConn struct {
	client    *http.Client
	req       *http.Request
	closeOnce sync.Once
	closed    chan struct{}
}

// httpConn is treated specially by Client.
func (hc *httpConn) LocalAddr() net.Addr              { return nullAddr }
func (hc *httpConn) RemoteAddr() net.Addr             { return nullAddr }
func (hc *httpConn) SetReadDeadline(time.Time) error  { return nil }
func (hc *httpConn) SetWriteDeadline(time.Time) error { return nil }
func (hc *httpConn) SetDeadline(time.Time) error      { return nil }
func (hc *httpConn) Write([]byte) (int, error)        { panic("Write called") }

func (hc *httpConn) Read(b []byte) (int, error) {
	<-hc.closed
	return 0, io.EOF
}

func (hc *httpConn) Close() error {
	hc.closeOnce.Do(func() { close(hc.closed) })
	return nil
}

// DefaultHTTPTimeouts represents the default timeout values used if further
// configuration is not provided.
var DefaultHTTPTimeouts = rpcTypes.HTTPTimeouts{
	ReadTimeout:  30 * time.Second,
	WriteTimeout: 30 * time.Second,
	IdleTimeout:  120 * time.Second,
}

// DialHTTPWithClient creates a new RPC client that connects to an RPC server over HTTP
// using the provided HTTP Client.
func DialHTTPWithClient(endpoint string, client *http.Client) (*Client, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", rpcTypes.ContentType)
	req.Header.Set("Accept", rpcTypes.ContentType)

	initctx := context.Background()
	return newClient(initctx, func(context.Context) (net.Conn, error) {
		return &httpConn{client: client, req: req, closed: make(chan struct{})}, nil
	})
}

// DialHTTP creates a new RPC client that connects to an RPC server over HTTP.
func DialHTTP(endpoint string) (*Client, error) {
	return DialHTTPWithClient(endpoint, new(http.Client))
}

func (c *Client) sendHTTP(ctx context.Context, op *requestOp, msg interface{}) error {
	hc := c.writeConn.(*httpConn)
	respBody, err := hc.doRequest(ctx, msg)
	if respBody != nil {
		defer respBody.Close()
	}

	if err != nil {
		if respBody != nil {
			buf := new(bytes.Buffer)
			if _, err2 := buf.ReadFrom(respBody); err2 == nil {
				return fmt.Errorf("%v %v", err, buf.String())
			}
		}
		return err
	}
	var respmsg rpcTypes.JsonrpcMessage
	if err := json.NewDecoder(respBody).Decode(&respmsg); err != nil {
		return err
	}
	op.resp <- &respmsg
	return nil
}

func (c *Client) sendBatchHTTP(ctx context.Context, op *requestOp, msgs []*rpcTypes.JsonrpcMessage) error {
	hc := c.writeConn.(*httpConn)
	respBody, err := hc.doRequest(ctx, msgs)
	if err != nil {
		return err
	}
	defer respBody.Close()
	var respmsgs []rpcTypes.JsonrpcMessage
	if err := json.NewDecoder(respBody).Decode(&respmsgs); err != nil {
		return err
	}
	for i := 0; i < len(respmsgs); i++ {
		op.resp <- &respmsgs[i]
	}
	return nil
}

func (hc *httpConn) doRequest(ctx context.Context, msg interface{}) (io.ReadCloser, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req := hc.req.WithContext(ctx)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.Body, errors.New(resp.Status)
	}
	return resp.Body, nil
}

// NewHTTPServer creates a new HTTP RPC server around an API provider.
//
// Deprecated: Server implements http.Handler
func NewHTTPServer(cors []string, vhosts []string, timeouts rpcTypes.HTTPTimeouts, srv *rpcTypes.Server) *http.Server {
	// Wrap the CORS-handler within a host-handler
	handler := newCorsHandler(srv, cors)
	handler = newVHostHandler(vhosts, handler)

	// Make sure timeout values are meaningful
	if timeouts.ReadTimeout < time.Second {
		log.Warn("Sanitizing invalid HTTP read timeout", "provided", timeouts.ReadTimeout, "updated", DefaultHTTPTimeouts.ReadTimeout)
		timeouts.ReadTimeout = DefaultHTTPTimeouts.ReadTimeout
	}
	if timeouts.WriteTimeout < time.Second {
		log.Warn("Sanitizing invalid HTTP write timeout", "provided", timeouts.WriteTimeout, "updated", DefaultHTTPTimeouts.WriteTimeout)
		timeouts.WriteTimeout = DefaultHTTPTimeouts.WriteTimeout
	}
	if timeouts.IdleTimeout < time.Second {
		log.Warn("Sanitizing invalid HTTP idle timeout", "provided", timeouts.IdleTimeout, "updated", DefaultHTTPTimeouts.IdleTimeout)
		timeouts.IdleTimeout = DefaultHTTPTimeouts.IdleTimeout
	}
	// Bundle and start the HTTP server
	return &http.Server{
		Handler:      handler,
		ReadTimeout:  timeouts.ReadTimeout,
		WriteTimeout: timeouts.WriteTimeout,
		IdleTimeout:  timeouts.IdleTimeout,
	}
}

func newCorsHandler(srv *rpcTypes.Server, allowedOrigins []string) http.Handler {
	// disable CORS support if user has not specified a custom CORS configuration
	if len(allowedOrigins) == 0 {
		return srv
	}
	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{http.MethodPost, http.MethodGet},
		MaxAge:         600,
		AllowedHeaders: []string{"*"},
	})
	return c.Handler(srv)
}

// virtualHostHandler is a handler which validates the Host-header of incoming requests.
// The virtualHostHandler can prevent DNS rebinding attacks, which do not utilize CORS-headers,
// since they do in-domain requests against the RPC api. Instead, we can see on the Host-header
// which domain was used, and validate that against a whitelist.
type virtualHostHandler struct {
	vhosts map[string]struct{}
	next   http.Handler
}

// ServeHTTP serves JSON-RPC requests over HTTP, implements http.Handler
func (h *virtualHostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if r.Host is not set, we can continue serving since a browser would set the Host header
	if r.Host == "" {
		h.next.ServeHTTP(w, r)
		return
	}
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		// Either invalid (too many colons) or no port specified
		host = r.Host
	}
	if ipAddr := net.ParseIP(host); ipAddr != nil {
		// It's an IP address, we can serve that
		h.next.ServeHTTP(w, r)
		return

	}
	// Not an ip address, but a hostname. Need to validate
	if _, exist := h.vhosts["*"]; exist {
		h.next.ServeHTTP(w, r)
		return
	}
	if _, exist := h.vhosts[host]; exist {
		h.next.ServeHTTP(w, r)
		return
	}
	http.Error(w, "invalid host specified", http.StatusForbidden)
}

func newVHostHandler(vhosts []string, next http.Handler) http.Handler {
	vhostMap := make(map[string]struct{})
	for _, allowedHost := range vhosts {
		vhostMap[strings.ToLower(allowedHost)] = struct{}{}
	}
	return &virtualHostHandler{vhostMap, next}
}
