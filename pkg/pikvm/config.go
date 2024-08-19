package pikvm

type PiKvmConfig struct {
	SkipVerify    bool   `env:"PIKVM_SKIP_TLS_VERIFY, default=false"`
	PiKvmHost     string `env:"PIKVM_HOST, required"`
	PiKvmUsername string `env:"PIKVM_USERNAME, default=admin"`
	PiKvmPassword string `env:"PIKVM_PASSWORD, default=admin"`
}
