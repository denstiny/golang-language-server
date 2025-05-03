package engine

import (
	"context"
	"fmt"
	"github.com/denstiny/golang-language-server/biz/flags"
	"log"
	"net"

	"github.com/sourcegraph/jsonrpc2"
)

type RouteFunc func(ctx context.Context, c *LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error)
type LspService struct {
	Config Config
	route  map[string]RouteFunc
}

func NewClient(r map[string]RouteFunc) *LspService {
	return &LspService{
		route: r,
	}
}

func (c *LspService) SetConfig(cfg Config) {
	c.Config = cfg
}

func (c *LspService) Start() {
	ctx := context.Background()
	if flags.SERVICE_TCP {
		c.TcpStart(ctx)
	} else if flags.SERVICE_STDIO {
		c.StdioStart(ctx)
	}
}

func (c *LspService) StdioStart(ctx context.Context) {
}

func (c *LspService) TcpStart(ctx context.Context) {
	h := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		log.Println("rpc msg:", req.Method)
		ctx = c.registerContext(ctx)
		result, err := c.Handle(ctx, conn, req)
		if err != nil {
			log.Fatal(err)
		}
		return result, err
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Config.ServerPort))
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

		defer conn.Close()
		go func(conn net.Conn) {
			rpcConn := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), h)
			if rpcConn == nil {
				return
			}
			log.Printf("Failed to start LSP server: %+v\n", rpcConn)
		}(conn)
	}
}

func (c *LspService) Register(method string, handler func(ctx context.Context, client *LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error)) {
	c.route[method] = handler
}

func (c *LspService) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	if handler, ok := c.route[req.Method]; ok {
		return handler(ctx, c, conn, req)
	}
	return nil, fmt.Errorf("method not found: %s", req.Method)
}

func (c *LspService) registerContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, "client", c)
	ctx = context.WithValue(ctx, "port", c.Config.ServerPort)
	return ctx
}
