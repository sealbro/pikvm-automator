package config

import (
	"github.com/sealbro/pikvm-automator/internal/grpc_ext"
	"github.com/sealbro/pikvm-automator/pkg/pikvm"
)

type PiKvmAutomatorConfig struct {
	pikvm.PiKvmConfig
	grpc_ext.GatewayConfig
	CommandsPath        string `env:"COMMANDS_PATH, default=commands.yaml"`
	TemplateMaxDeep     int    `env:"TEMPLATE_MAX_DEEP, default=10"`
	CallDebounceSeconds int    `env:"CALL_DEBOUNCE_SECONDS, default=2"`
	DatabasePath        string `env:"DATABASE_PATH, default=pikvm-automator.db"`
}
