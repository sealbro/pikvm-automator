package server

import (
	"context"
	gen "github.com/sealbro/pikvm-automator/generated/go"
	"github.com/sealbro/pikvm-automator/internal/macro"
	"github.com/sealbro/pikvm-automator/internal/queue"
	"github.com/sealbro/pikvm-automator/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type PiKvmAutomatorServer struct {
	gen.UnimplementedPiKvmAutomatorServer
	player            *queue.ExpressionPlayer
	commandRepository *storage.CommandRepository
}

func NewPiKvmAutomatorServer(player *queue.ExpressionPlayer, commandRepository *storage.CommandRepository) *PiKvmAutomatorServer {
	return &PiKvmAutomatorServer{
		player:            player,
		commandRepository: commandRepository,
	}
}

func (s *PiKvmAutomatorServer) CommandList(context.Context, *gen.CommandListRequest) (*gen.CommandListResponse, error) {
	commands, err := s.commandRepository.GetCommands()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get commands")
	}

	resp := &gen.CommandListResponse{}
	for _, cmd := range commands {
		resp.Commands = append(resp.Commands, &gen.Command{
			Id:          cmd.ID,
			Description: cmd.Description,
			Expression:  cmd.Expression,
		})
	}

	return resp, nil
}
func (s *PiKvmAutomatorServer) CallCommand(_ context.Context, req *gen.CallCommandRequest) (*gen.CallCommandResponse, error) {
	if req.Expression == "" {
		return nil, status.Errorf(codes.InvalidArgument, "expression is required")
	}

	expressions := req.Expression

	compile, err := regexp.Compile("%([\\w\\s_]+)%")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't compile regexp")
	}
	matches := compile.FindAllStringSubmatch(expressions, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			command, err := s.commandRepository.GetCommand(match[1])
			if err != nil {
				expressions = compile.ReplaceAllString(expressions, "42ms")
			} else {
				expressions = compile.ReplaceAllString(expressions, command.Expression)
			}
		}
	}

	macroExp := macro.New(expressions)
	s.player.AddExpression(macroExp)

	return &gen.CallCommandResponse{}, nil
}

func (s *PiKvmAutomatorServer) DeleteCommand(_ context.Context, req *gen.DeleteCommandRequest) (*gen.DeleteCommandResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "command is required")
	}

	err := s.commandRepository.DeleteCommand(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "command not found")
	}

	return &gen.DeleteCommandResponse{}, nil
}
func (s *PiKvmAutomatorServer) CreateCommand(_ context.Context, req *gen.CreateCommandRequest) (*gen.CreateCommandResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "command is required")
	}

	command, _ := s.commandRepository.GetCommand(req.Id)
	if command != nil {
		return nil, status.Errorf(codes.AlreadyExists, "command already exists")
	}

	err := s.commandRepository.CreateCommand(storage.Command{
		ID:          req.Id,
		Description: req.Description,
		Expression:  req.Expression,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't create command")
	}

	return &gen.CreateCommandResponse{}, nil
}
