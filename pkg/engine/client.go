package engine

import (
	"context"
	"fmt"
	"github.com/denstiny/golang-language-server/biz/flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"os"

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
		log.Info().Msg("handle msg:" + req.Method)
		ctx = c.registerContext(ctx)
		ctx = context.WithValue(ctx, rpc_conn, conn)
		result, err := c.Handle(ctx, conn, req)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		return result, err
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Config.ServerPort))
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Msg(err.Error())
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

const rpc_conn = "rpc-client"

func GetRpcConn(ctx context.Context) *jsonrpc2.Conn {
	if v := ctx.Value(rpc_conn); v != nil {
		if conn, ok := v.(*jsonrpc2.Conn); ok {
			return conn
		}
	}
	return nil
}

func init() {
	// 设置彩色输出
	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}
	log.Logger = log.Output(output).Level(zerolog.DebugLevel)
}
