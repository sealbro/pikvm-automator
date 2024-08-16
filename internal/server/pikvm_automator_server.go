package server

import (
	"context"
	gen "github.com/sealbro/pikvm-automator/generated/go"
	"github.com/sealbro/pikvm-automator/internal/macro"
	"github.com/sealbro/pikvm-automator/internal/queue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PiKvmAutomatorServer struct {
	gen.UnimplementedPiKvmAutomatorServer
	player *queue.ExpressionPlayer
}

func NewPiKvmAutomatorServer(player *queue.ExpressionPlayer) *PiKvmAutomatorServer {
	return &PiKvmAutomatorServer{
		player: player,
	}
}

func (s *PiKvmAutomatorServer) CommandList(context.Context, *gen.CommandListRequest) (*gen.CommandListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandList not implemented")
}
func (s *PiKvmAutomatorServer) CallCommand(_ context.Context, req *gen.CallCommandRequest) (*gen.CallCommandResponse, error) {
	if req.Expression == "" {
		return nil, status.Errorf(codes.InvalidArgument, "expression is required")
	}

	expression := macro.New(req.Expression)
	s.player.AddExpression(expression)

	return &gen.CallCommandResponse{}, nil
}

func (s *PiKvmAutomatorServer) DeleteCommand(context.Context, *gen.DeleteCommandRequest) (*gen.DeleteCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCommand not implemented")
}
func (s *PiKvmAutomatorServer) CreateCommand(context.Context, *gen.CreateCommandRequest) (*gen.CreateCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCommand not implemented")
}
