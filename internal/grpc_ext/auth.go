package grpc_ext

import (
	"context"
	"encoding/base64"
	"github.com/sealbro/pikvm-automator/pkg/pikvm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"strings"
	"time"
)

const (
	authorizationHeader = "authorization"
)

type AuthInterceptor struct {
	piKvmClient *pikvm.PiKvmClient
	logger      *slog.Logger
	tokenExpire map[string]time.Time
}

func NewAuthInterceptor(logger *slog.Logger, piKvmClient *pikvm.PiKvmClient) *AuthInterceptor {
	return &AuthInterceptor{
		tokenExpire: make(map[string]time.Time, 1),
		logger:      logger,
		piKvmClient: piKvmClient,
	}
}

func (a *AuthInterceptor) Interceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	// Check if the token is valid
	if !a.valid(ctx, md[authorizationHeader]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// Continue execution of handler after successful authentication
	return handler(ctx, req)
}

func (a *AuthInterceptor) valid(ctx context.Context, tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	authToken := tokens[0]
	expireAfter := time.Hour
	if v, ok := a.tokenExpire[authToken]; ok && time.Since(v) < expireAfter {
		return true
	}

	isBasic := strings.HasPrefix(authToken, "Basic ")
	if !isBasic {
		a.logger.WarnContext(ctx, "invalid token", slog.Any("token", authToken))
		return false
	}

	base64Credentials := strings.TrimSpace(strings.TrimPrefix(authToken, "Basic "))
	decodeString, err := base64.StdEncoding.DecodeString(base64Credentials)
	if err != nil {
		a.logger.WarnContext(ctx, "decode base64", slog.Any("err", err))
		return false
	}

	credentials := strings.Split(string(decodeString), ":")
	if len(credentials) != 2 {
		a.logger.WarnContext(ctx, "invalid credentials", slog.Any("credentials", credentials))
		return false
	}

	username := credentials[0]
	password := credentials[1]

	check := a.piKvmClient.Check(ctx, username, password)
	if check {
		a.tokenExpire[authToken] = time.Now()
	}

	return check
}
