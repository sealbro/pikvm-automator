package queue

import (
	"context"
	"github.com/sealbro/pikvm-automator/internal/macro"
	"github.com/sealbro/pikvm-automator/pkg/pikvm"
	"log/slog"
	"sync"
	"time"
)

type ExpressionPlayer struct {
	l           sync.Mutex
	Expressions []*macro.Expression
	logger      *slog.Logger
}

func NewExpressionPlayer(logger *slog.Logger) *ExpressionPlayer {
	return &ExpressionPlayer{
		logger: logger,
	}
}

func (p *ExpressionPlayer) AddExpression(expression *macro.Expression) {
	p.l.Lock()
	defer p.l.Unlock()
	p.Expressions = append(p.Expressions, expression)
}

func (p *ExpressionPlayer) Play(ctx context.Context) <-chan pikvm.PiKvmEvent {
	events := make(chan pikvm.PiKvmEvent, 20)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(events)
				return
			default:
				p.l.Lock()
				if len(p.Expressions) > 0 {
					exp := p.Expressions[0]
					p.Expressions = p.Expressions[1:]
					p.l.Unlock()
					group := exp.Parse()
					p.logger.InfoContext(ctx, "Expression run", slog.String("expression", exp.String()), slog.Duration("delay", group.TotalDelay()))
					play(group.Events, events)
					p.logger.InfoContext(ctx, "Expression played", slog.String("expression", exp.String()))
				} else {
					p.l.Unlock()
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	return events
}

func play(macros []macro.Macro, events chan<- pikvm.PiKvmEvent) {
	for _, m := range macros {
		switch v := m.(type) {
		case macro.Delay:
			time.Sleep(v.Time)
		case macro.KeyEvent:
			events <- pikvm.PiKvmEvent{
				EventType: pikvm.Keyboard,
				Event:     pikvm.KeyboardEvent{Key: v.Key, State: v.State},
			}
		case macro.MouseEvent:
			events <- pikvm.PiKvmEvent{
				EventType: pikvm.Keyboard,
				Event: pikvm.MouseEvent{
					To: pikvm.MousePoint{
						X: v.X,
						Y: v.Y,
					},
				},
			}
		case macro.Repeat:
			for i := 0; i < v.Repeats; i++ {
				play(v.Events, events)
			}
		case macro.Bind:
			play(v.Events, events)
		}
	}
}
