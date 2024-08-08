package pikvm

type PiKvmConfig struct {
	PiKvmHost     string `env:"PIKVM_HOST, required"`
	PiKvmUsername string `env:"PIKVM_USERNAME, default=admin"`
	PiKvmPassword string `env:"PIKVM_PASSWORD, default=admin"`
}
