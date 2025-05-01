package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/denstiny/golang-language-server/biz/conts"
	"github.com/sourcegraph/go-lsp"
	"log"
	"net"

	"github.com/sourcegraph/jsonrpc2"
)

type RouteFunc func(ctx context.Context, c *Client, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error)
type Client struct {
	conn   *jsonrpc2.Conn
	Config Config
	route  map[string]RouteFunc
}

type Config struct {
	Port int
}

func NewClient(r map[string]RouteFunc) *Client {
	return &Client{
		route: r,
	}
}

func (c *Client) Start() {
	ctx := context.Background()
	h := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		if req.Method != "initialize" {
			return c.Initialize(ctx, conn, req)
		} else {
			ctx = c.registerContext(ctx)
			result, err := c.Handle(ctx, conn, req)
			if err != nil {
				log.Fatalf("failed to handle request: %v", err)
				return lsp.None{}, nil
			}
			return result, nil
		}
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Config.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		c.conn = jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), h)
		go func() {
			if err := c.conn; err != nil {
				log.Fatalf("Failed to start LSP server: %v", err)
			}
		}()
	}
}

func (c *Client) Register(method string, handler func(ctx context.Context, client *Client, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error)) {
	c.route[method] = handler
}

func (c *Client) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	if handler, ok := c.route[req.Method]; ok {
		return handler(ctx, c, conn, req)
	}
	return nil, fmt.Errorf("method not found: %s", req.Method)
}

func (c *Client) registerContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, "client", c)
	ctx = context.WithValue(ctx, "port", c.Config.Port)
	return ctx
}

func (c *Client) Initialize(ctx context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	var params lsp.InitializeParams
	err := json.Unmarshal(*req.Params, &params)
	if err != nil {
		return nil, err
	}
	log.Println("Initialize params:", params)

	return lsp.InitializeResult{Capabilities: conts.ServerCapabilities}, nil
}
