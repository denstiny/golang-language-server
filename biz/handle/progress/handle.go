package progress

import (
	"context"
	"fmt"
	"github.com/denstiny/golang-language-server/pkg/engine"
	"pkg.nimblebun.works/go-lsp"
	"sync"
)

var mu sync.Mutex

func notify(ctx context.Context, param lsp.ProgressParams) error {
	mu.Lock()
	defer mu.Unlock()
	conn := engine.GetRpcConn(ctx)
	if conn == nil {
		return fmt.Errorf("notify fail: rpc conn is nil")
	}

	err := conn.Notify(ctx, "$/progress", param)
	if err != nil {
		return err
	}
	return nil
}

type Progress struct {
	Token lsp.ProgressToken `json:"token"`
	Title string            `json:"title"`
}

func NewProgress(token lsp.ProgressToken, title string) *Progress {
	return &Progress{
		Token: token,
		Title: title,
	}
}

func (p *Progress) Begin(ctx context.Context, message string, cancellable bool) error {
	return notify(ctx, lsp.ProgressParams{
		Token: p.Token,
		Value: lsp.WorkDoneProgressBegin{
			Title:       p.Title,
			Message:     message,
			Cancellable: cancellable,
			Percentage:  0,
		},
	})
}

func (p *Progress) End(ctx context.Context, kind string, message string) error {
	return notify(ctx, lsp.ProgressParams{
		Token: p.Token,
		Value: lsp.WorkDoneProgressEnd{
			Kind:    kind,
			Message: message,
		},
	})
}

func (p *Progress) Update(ctx context.Context, kind string, message string, percentage int) error {
	return notify(ctx, lsp.ProgressParams{
		Token: p.Token,
		Value: lsp.WorkDoneProgressReport{
			Kind:        kind,
			Message:     message,
			Percentage:  percentage,
			Cancellable: false,
		},
	})
}
