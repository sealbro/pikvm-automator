package config

import (
	"github.com/sealbro/pikvm-automator/internal/grpc_ext"
	"github.com/sealbro/pikvm-automator/pkg/pikvm"
)

type PiKvmAutomatorConfig struct {
	pikvm.PiKvmConfig
	grpc_ext.GatewayConfig
}
