package server

import (
	"context"
	gen "github.com/sealbro/pikvm-automator/generated/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PiKvmAutomatorServer struct {
	gen.UnimplementedPiKvmAutomatorServer
}

func NewPiKvmAutomatorServer() *PiKvmAutomatorServer {
	return &PiKvmAutomatorServer{}
}

func (s *PiKvmAutomatorServer) CommandList(context.Context, *gen.CommandListRequest) (*gen.CommandListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandList not implemented")
}
func (s *PiKvmAutomatorServer) CallCommand(context.Context, *gen.CallCommandRequest) (*gen.CallCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallCommand not implemented")
}
func (s *PiKvmAutomatorServer) ManualCall(context.Context, *gen.ManualCallRequest) (*gen.ManualCallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ManualCall not implemented")
}
func (s *PiKvmAutomatorServer) DeleteCommand(context.Context, *gen.DeleteCommandRequest) (*gen.DeleteCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCommand not implemented")
}
func (s *PiKvmAutomatorServer) CreateCommand(context.Context, *gen.CreateCommandRequest) (*gen.CreateCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCommand not implemented")
}
func (s *PiKvmAutomatorServer) UpdateCommand(context.Context, *gen.UpdateCommandRequest) (*gen.UpdateCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCommand not implemented")
}
