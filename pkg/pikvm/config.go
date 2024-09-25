package pikvm

type PiKvmConfig struct {
	SkipVerify    bool   `env:"PIKVM_SKIP_TLS_VERIFY, default=false"`
	PiKvmAddress  string `env:"PIKVM_ADDRESS, required"`
	PiKvmUsername string `env:"PIKVM_USERNAME, default=admin"`
	PiKvmPassword string `env:"PIKVM_PASSWORD, default=admin"`
}
