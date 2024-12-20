package pikvm

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sealbro/pikvm-automator/pkg/rand"
	"log/slog"
	"net"
	"net/http"
	"syscall"
	"time"
)

// PiKvmClient is a client for PiKVM see expected events https://docs.pikvm.org/api/#websocket-events
type PiKvmClient struct {
	config     PiKvmConfig
	logger     *slog.Logger
	connection *websocket.Conn
	rnd        *rand.Rand
}

func NewPiKvmClient(logger *slog.Logger, config PiKvmConfig) *PiKvmClient {
	return &PiKvmClient{
		config: config,
		logger: logger,
		rnd:    rand.New(),
	}
}

func (c *PiKvmClient) Check(ctx context.Context, username, password string) bool {
	timeout, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()

	address := c.config.ApiAddress("/api/auth/check")
	request, err := http.NewRequestWithContext(timeout, http.MethodGet, address, nil)
	if err != nil {
		c.logger.Error("request check", slog.Any("err", err))
		return false
	}

	request.Header = piKvmAuthHeader(username, password)

	do, err := http.DefaultClient.Do(request)
	if err != nil {
		c.logger.Error("do request check", slog.Any("err", err))
		return false
	}

	return do.StatusCode == http.StatusOK
}

func (c *PiKvmClient) StartWebSocket(ctx context.Context, sender <-chan PiKvmEvent, receiver func([]byte)) error {
	err := c.reconnect()
	if err != nil {
		return err
	}

	go c.receiveEvent(ctx, receiver)
	go c.sendEvent(ctx, sender)
	go c.stop(ctx)

	return nil
}

func (c *PiKvmClient) stop(ctx context.Context) {
	<-ctx.Done()
	if c.connection != nil {
		err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			c.logger.ErrorContext(ctx, "stop", slog.Any("err", err))
			return
		}
		_ = c.connection.Close()
		c.connection = nil
	}
}

func (c *PiKvmClient) receiveEvent(ctx context.Context, receive func([]byte)) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, message, err := c.connection.ReadMessage()
			if err != nil {
				c.logger.ErrorContext(ctx, "receiveEvent", slog.Any("err", err))
				return
			}
			receive(message)
		}
	}
}

func (c *PiKvmClient) sendEvent(ctx context.Context, sender <-chan PiKvmEvent) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-sender:
			data, err := json.Marshal(event)
			if err != nil {
				c.logger.ErrorContext(ctx, "marshal", slog.Any("err", err))
				continue
			}

			// random delay
			time.Sleep(time.Duration(c.rnd.Range(1, 5)) * time.Millisecond)

			err = c.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				if errors.Is(err, syscall.EPIPE) {
					c.logger.WarnContext(ctx, "reconnecting...", slog.Any("err", err))
					err = c.reconnect()
					if err != nil {
						c.logger.ErrorContext(ctx, "reconnect failed", slog.Any("err", err))
						return
					}

					err = c.connection.WriteMessage(websocket.TextMessage, data)
					if err != nil {
						c.logger.ErrorContext(ctx, "reconnect write", slog.Any("err", err))
						return
					}
				}

				c.logger.ErrorContext(ctx, "sendEvent", slog.Any("err", err))
				continue
			} else {
				c.logger.InfoContext(ctx, "sendEvent", slog.String("data", string(data)))
			}
		}
	}
}

func (c *PiKvmClient) reconnect() error {
	if c.connection != nil {
		_ = c.connection.Close()
	}

	c.logger.Info("connecting to", slog.String("url", c.config.PiKvmAddress))

	conn, err := c.dial()
	if err != nil {
		return fmt.Errorf("pikvm dial: %w", err)
	}
	keepAlive(conn, 30*time.Second)

	c.logger.Info("connected to", slog.String("url", c.config.PiKvmAddress))
	c.connection = conn

	return nil
}

func (c *PiKvmClient) dial() (wsConn *websocket.Conn, err error) {
	authHeader := c.config.AuthHeader()
	if c.config.IsUnixSocket() {
		sockPath, params, wsAddress := c.config.UnixSocketAddress(0)

		c.logger.Info("connecting to", slog.String("url", sockPath), slog.String("params", params))

		d := websocket.Dialer{
			NetDialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", sockPath)
			},
		}
		wsConn, _, err = d.Dial(wsAddress, authHeader)
	} else {
		websocket.DefaultDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: c.config.SkipVerify}
		wsConn, _, err = websocket.DefaultDialer.Dial(c.config.WebSocketAddress(0), authHeader)
	}

	if err != nil {
		return nil, err
	}
	return wsConn, nil
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		for {
			err := c.WriteMessage(websocket.PingMessage, []byte("{\"event_type\": \"ping\", \"event\": {}}"))
			if err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Since(lastResponse) > timeout {
				_ = c.Close()
				return
			}
		}
	}()
}
