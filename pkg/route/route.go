package route

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
)

type Route struct {
	route map[string]func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error)
}

func NewRoute() *Route {
	return &Route{
		route: make(map[string]func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error)),
	}
}

func (r *Route) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	if handler, ok := r.route[req.Method]; ok {
		return handler(ctx, conn, req)
	}
	return nil, nil
}

func (r *Route) Register(method string, handler func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error)) {
	r.route[method] = handler
}
