package server

import (
	"context"
	gen "github.com/sealbro/pikvm-automator/generated/go"
	"github.com/sealbro/pikvm-automator/internal/macro"
	"github.com/sealbro/pikvm-automator/internal/queue"
	"github.com/sealbro/pikvm-automator/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

var (
	templateRegexp, _ = regexp.Compile("%([\\w\\s_]+)%")
)

type PiKvmAutomatorServer struct {
	gen.UnimplementedPiKvmAutomatorServer
	player            *queue.ExpressionPlayer
	commandRepository *storage.CommandRepository
	logger            *slog.Logger
	lastCall          time.Time
}

func NewPiKvmAutomatorServer(logger *slog.Logger, player *queue.ExpressionPlayer, commandRepository *storage.CommandRepository) *PiKvmAutomatorServer {
	return &PiKvmAutomatorServer{
		logger:            logger,
		player:            player,
		commandRepository: commandRepository,
		lastCall:          time.Now(),
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
func (s *PiKvmAutomatorServer) CallCommand(ctx context.Context, req *gen.CallCommandRequest) (*gen.CallCommandResponse, error) {
	if req.Expression == "" {
		return nil, status.Errorf(codes.InvalidArgument, "expression is required")
	}

	if time.Since(s.lastCall) < 1*time.Second {
		s.logger.WarnContext(ctx, "too many requests", slog.String("expression", req.Expression))
		return nil, status.Errorf(codes.ResourceExhausted, "too many requests")
	}
	s.lastCall = time.Now()

	expressions := req.Expression
	maxDeep := 5
	for {
		matches := templateRegexp.FindAllStringSubmatch(expressions, -1)
		for _, match := range matches {
			template := match[0]
			id := match[1]
			command, err := s.commandRepository.GetCommand(id)
			if err != nil {
				s.logger.WarnContext(ctx, "can't get command", slog.String("id", id), slog.Any("err", err))
				expressions = strings.ReplaceAll(expressions, template, "42ms")
			} else {
				expressions = strings.ReplaceAll(expressions, template, command.Expression)
			}
		}

		if len(matches) == 0 {
			break
		}

		if maxDeep == 0 {
			s.logger.WarnContext(ctx, "max deep reached", slog.String("expression", expressions), slog.Int("maxDeep", maxDeep))
			expressions = templateRegexp.ReplaceAllString(expressions, "42ms")
			break
		}

		maxDeep--
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
