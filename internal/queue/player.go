package queue

import (
	"context"
	"github.com/sealbro/pikvm-automator/internal/macro"
	"github.com/sealbro/pikvm-automator/pkg/pikvm"
	"github.com/sealbro/pikvm-automator/pkg/pikvm/screen"
	"log/slog"
	"sync"
	"time"
)

type ExpressionPlayer struct {
	l           sync.Mutex
	Expressions []*macro.Expression
	logger      *slog.Logger
	screen      *screen.Position
}

func NewExpressionPlayer(logger *slog.Logger) *ExpressionPlayer {
	return &ExpressionPlayer{
		screen: screen.NewFullHD(),
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
					p.play(group.Events, events)
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

func (p *ExpressionPlayer) play(macros []macro.Macro, events chan<- pikvm.PiKvmEvent) {
	for _, m := range macros {
		switch v := m.(type) {
		case macro.Delay:
			time.Sleep(v.Time)
		case macro.KeyPressEvent:
			events <- pikvm.PiKvmEvent{
				EventType: pikvm.Keyboard,
				Event:     pikvm.KeyboardEvent{Key: v.Key, State: v.State},
			}
		case macro.MouseMoveEvent:
			x, y := p.screen.ToPiKvmPoints(v.X, v.Y)
			events <- pikvm.PiKvmEvent{
				EventType: pikvm.MouseMove,
				Event: pikvm.MouseMoveEvent{
					To: pikvm.MousePoint{
						X: x,
						Y: y,
					},
				},
			}
		case macro.MouseClickEvent:
			events <- pikvm.PiKvmEvent{
				EventType: pikvm.MouseButton,
				Event:     pikvm.MouseButtonEvent{Button: v.Button, State: v.State},
			}
		case macro.Repeat:
			for i := 0; i < v.Repeats; i++ {
				p.play(v.Events, events)
			}
		case macro.Bind:
			p.play(v.Events, events)
		}
	}
}
