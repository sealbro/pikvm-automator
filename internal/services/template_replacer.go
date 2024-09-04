package services

import (
	"context"
	"github.com/sealbro/pikvm-automator/internal/config"
	"github.com/sealbro/pikvm-automator/internal/macro"
	"github.com/sealbro/pikvm-automator/internal/storage"
	"log/slog"
	"regexp"
	"strings"
)

var (
	templateRegexp, _ = regexp.Compile("%([\\w\\s_]+)%")
)

type TemplateReplacer struct {
	maxDeep           int
	logger            *slog.Logger
	commandRepository *storage.CommandRepository
}

func NewTemplateReplacer(logger *slog.Logger, commandRepository *storage.CommandRepository, config config.PiKvmAutomatorConfig) *TemplateReplacer {
	return &TemplateReplacer{
		logger:            logger,
		commandRepository: commandRepository,
		maxDeep:           config.TemplateMaxDeep,
	}
}

func (t *TemplateReplacer) Replace(ctx context.Context, expressions string) *macro.Expression {
	maxDeep := t.maxDeep
	for {
		matches := templateRegexp.FindAllStringSubmatch(expressions, -1)
		for _, match := range matches {
			template := match[0]
			id := match[1]
			command, err := t.commandRepository.GetCommand(id)
			if err != nil {
				t.logger.WarnContext(ctx, "can't get command", slog.String("id", id), slog.Any("err", err))
				expressions = strings.ReplaceAll(expressions, template, "42ms")
			} else {
				expressions = strings.ReplaceAll(expressions, template, command.Expression)
			}
		}

		if len(matches) == 0 {
			break
		}

		if maxDeep == 0 {
			t.logger.WarnContext(ctx, "max deep reached", slog.String("expression", expressions), slog.Int("maxDeep", maxDeep))
			expressions = templateRegexp.ReplaceAllString(expressions, "42ms")
			break
		}

		maxDeep--
	}

	return macro.NewExpression(expressions)
}
