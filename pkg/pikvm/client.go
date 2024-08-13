package pikvm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sealbro/pikvm-automator/pkg/rand"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

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

func (c *PiKvmClient) Start(ctx context.Context, send <-chan PiKvmEvent) (error, <-chan []byte) {
	u := url.URL{Scheme: "wss", Host: c.config.PiKvmHost, Path: "/api/ws", RawQuery: "stream=0"}
	c.logger.Info("connecting to ", slog.String("url", u.String()))

	httpHeader := http.Header{}
	httpHeader.Add("X-KVMD-User", c.config.PiKvmUsername)
	httpHeader.Add("X-KVMD-Passwd", c.config.PiKvmPassword)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), httpHeader)
	if err != nil {
		return fmt.Errorf("pikvm dial: %w", err), nil
	}

	c.connection = conn

	receive := make(chan []byte)
	go c.receiveEvent(ctx, receive)
	go c.sendEvent(ctx, send)
	go c.stop(ctx)

	return nil, receive
}

func (c *PiKvmClient) stop(ctx context.Context) error {
	<-ctx.Done()
	if c.connection != nil {
		err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			return fmt.Errorf("pikvm write close: %w", err)
		}
		_ = c.connection.Close()
		c.connection = nil
		return nil
	}
	return nil
}

func (c *PiKvmClient) receiveEvent(ctx context.Context, receive chan<- []byte) {
	defer close(receive)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, message, err := c.connection.ReadMessage()
			if err != nil {
				c.logger.ErrorContext(ctx, "read", slog.Any("err", err))
				return
			}
			receive <- message
		}
	}
}

func (c *PiKvmClient) sendEvent(ctx context.Context, send <-chan PiKvmEvent) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-send:
			data, err := json.Marshal(event)
			if err != nil {
				c.logger.ErrorContext(ctx, "marshal", slog.Any("err", err))
				continue
			}

			// random delay
			time.Sleep(time.Duration(c.rnd.Range(1, 5)) * time.Millisecond)

			err = c.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				c.logger.ErrorContext(ctx, "write", slog.Any("err", err))
				continue
			} else {
				c.logger.InfoContext(ctx, "send", slog.String("data", string(data)))
			}
		}
	}
}
