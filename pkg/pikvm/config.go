package pikvm

import (
	"fmt"
	"net/http"
	"net/url"
)

type PiKvmConfig struct {
	SkipVerify    bool   `env:"PIKVM_SKIP_TLS_VERIFY, default=false"`
	PiKvmAddress  string `env:"PIKVM_ADDRESS, required"`
	PiKvmSource   string `env:"PIKVM_SOURCE, default=wss"`
	PiKvmUsername string `env:"PIKVM_USERNAME, default=admin"`
	PiKvmPassword string `env:"PIKVM_PASSWORD, default=admin"`
}

// ApiAddress returns the http address for the api
func (c PiKvmConfig) ApiAddress(apiPath string) string {
	result, _ := url.JoinPath("https://", c.PiKvmAddress, apiPath)
	return result
}

// WebSocketAddress returns the address for the websocket connection
// example: wss://<IP or HOST>/api/ws?stream=0
func (c PiKvmConfig) WebSocketAddress(stream int) string {
	return fmt.Sprintf("wss://%s/api%s", c.PiKvmAddress, streamParams(stream))
}

// UnixSocketAddress returns the address for the unix socket connection
// example: ws://unix:/run/kvmd/kvmd.sock/ws?stream=0
func (c PiKvmConfig) UnixSocketAddress(stream int) (sockPath string, queryParams string, wsAddress string) {
	queryParams = streamParams(stream)
	sockPath = c.PiKvmAddress
	wsAddress = fmt.Sprintf("ws://unix%s", queryParams)

	return sockPath, queryParams, wsAddress
}

// IsUnixSocket returns true if the source is unix
func (c PiKvmConfig) IsUnixSocket() bool {
	return c.PiKvmSource == "unix"
}

// AuthHeader returns the auth header for the connection
func (c PiKvmConfig) AuthHeader() http.Header {
	return piKvmAuthHeader(c.PiKvmUsername, c.PiKvmPassword)
}

func streamParams(stream int) string {
	return fmt.Sprintf("/ws?stream=%d", stream)
}

func piKvmAuthHeader(username, password string) http.Header {
	httpHeader := http.Header{}
	httpHeader.Add("X-KVMD-User", username)
	httpHeader.Add("X-KVMD-Passwd", password)

	return httpHeader
}
