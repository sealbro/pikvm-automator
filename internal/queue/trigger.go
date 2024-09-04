package queue

import (
	"github.com/sealbro/pikvm-automator/internal/macro"
	"log/slog"
	"slices"
	"sync"
	"time"
)

type TriggerType string

const (
	PiKvmHidOnline TriggerType = "pikvm_hid_online"
)

var (
	triggers = []TriggerType{
		PiKvmHidOnline,
	}

	triggerTypeTimeouts = map[TriggerType]time.Duration{
		PiKvmHidOnline: 10 * time.Minute,
	}
)

type triggerItem struct {
	expressions *macro.Expression
	expire      time.Time
}

type ExpressionTrigger struct {
	l        sync.Mutex
	triggers map[TriggerType]triggerItem
	logger   *slog.Logger
	player   *ExpressionPlayer
}

func NewExpressionTrigger(logger *slog.Logger, player *ExpressionPlayer) *ExpressionTrigger {
	return &ExpressionTrigger{
		player:   player,
		logger:   logger,
		triggers: make(map[TriggerType]triggerItem),
	}
}

func (t *ExpressionTrigger) Rise(triggerType TriggerType) {
	t.l.Lock()
	defer t.l.Unlock()
	if item, ok := t.triggers[triggerType]; ok {
		if time.Now().Before(item.expire) {
			t.player.AddExpression(item.expressions)
		} else {
			t.logger.Warn("Trigger expired", slog.String("trigger", string(triggerType)))
		}
		delete(t.triggers, triggerType)
	}
}

func (t *ExpressionTrigger) AddExpression(triggerType TriggerType, expression *macro.Expression) {
	t.l.Lock()
	defer t.l.Unlock()
	t.logger.Info("Added trigger", slog.String("trigger", string(triggerType)))
	t.triggers[triggerType] = triggerItem{
		expressions: expression,
		expire:      time.Now().Add(triggerTypeTimeouts[triggerType]),
	}
}

func (t TriggerType) IsValid() bool {
	return slices.Contains(triggers, t)
}
