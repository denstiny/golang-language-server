package engine

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/denstiny/golang-language-server/pkg/route"
	"github.com/sourcegraph/jsonrpc2"
)

type Client struct {
	route  *route.Route
	conn   *jsonrpc2.Conn
	Config Config
}

type Config struct {
	Port int
}

func NewClient(r *route.Route) *Client {
	return &Client{
		route: r,
	}
}

func (c *Client) Start() {
	ctx := context.Background()
	h := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		return c.route.Handle(ctx, conn, req)
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
