package server

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	gen "github.com/sealbro/pikvm-automator/generated/go"
	"github.com/sealbro/pikvm-automator/internal/config"
	"github.com/sealbro/pikvm-automator/internal/queue"
	"github.com/sealbro/pikvm-automator/internal/repository"
	"github.com/sealbro/pikvm-automator/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

type PiKvmAutomatorServer struct {
	gen.UnimplementedPiKvmAutomatorServer
	player            *queue.ExpressionPlayer
	commandRepository *repository.Queries
	templateReplacer  *services.TemplateReplacer
	logger            *slog.Logger
	lastCall          time.Time
	callDebounce      time.Duration
	trigger           *queue.ExpressionTrigger
}

func NewPiKvmAutomatorServer(logger *slog.Logger, player *queue.ExpressionPlayer, commandRepository *repository.Queries, templateReplacer *services.TemplateReplacer, trigger *queue.ExpressionTrigger, config config.PiKvmAutomatorConfig) *PiKvmAutomatorServer {
	return &PiKvmAutomatorServer{
		logger:            logger,
		player:            player,
		commandRepository: commandRepository,
		lastCall:          time.Now(),
		callDebounce:      time.Duration(config.CallDebounceSeconds) * time.Second,
		templateReplacer:  templateReplacer,
		trigger:           trigger,
	}
}

func (s *PiKvmAutomatorServer) CommandList(ctx context.Context, _ *gen.CommandListRequest) (*gen.CommandListResponse, error) {
	commands, err := s.commandRepository.ListCommands(ctx)
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

	if time.Since(s.lastCall) < s.callDebounce {
		s.logger.WarnContext(ctx, "too many requests", slog.String("expression", req.Expression))
		return nil, status.Errorf(codes.ResourceExhausted, "too many requests")
	}
	s.lastCall = time.Now()

	macroExp := s.templateReplacer.Replace(ctx, req.Expression)
	if req.Trigger != "" {
		if queue.TriggerType(req.Trigger).IsValid() {
			s.trigger.AddExpression(queue.TriggerType(req.Trigger), macroExp)
		} else {
			s.logger.WarnContext(ctx, "Invalid trigger", slog.String("trigger", req.Trigger))
		}
	} else {
		s.player.AddExpression(macroExp)
	}

	return &gen.CallCommandResponse{}, nil
}

func (s *PiKvmAutomatorServer) DeleteCommand(ctx context.Context, req *gen.DeleteCommandRequest) (*gen.DeleteCommandResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "command is required")
	}

	err := s.commandRepository.DeleteCommand(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "command not found")
	}

	return &gen.DeleteCommandResponse{}, nil
}
func (s *PiKvmAutomatorServer) CreateCommand(ctx context.Context, req *gen.CreateCommandRequest) (*gen.CreateCommandResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "command is required")
	}

	_, err := s.commandRepository.GetCommand(ctx, req.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("can't get command", slog.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, "can't get command")
	}

	_, err = s.commandRepository.CreateCommand(ctx, repository.CreateCommandParams{
		ID:          req.Id,
		Description: req.Description,
		Expression:  req.Expression,
	})
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
			return nil, status.Errorf(codes.AlreadyExists, "command already exists")
		}

		s.logger.Error("can't create command", slog.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, "can't create command")
	}

	return &gen.CreateCommandResponse{}, nil
}
