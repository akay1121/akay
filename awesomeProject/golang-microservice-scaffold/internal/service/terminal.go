package service

import (
	"context"
	v1 "example/api/terminal_server"
	"example/internal/biz"
	"time"
)

type TerminalService struct {
	v1.UnimplementedTerminalManagementServer
	mgr *biz.TerminalManager
}

func NewTerminalService(mgr *biz.TerminalManager) *TerminalService {
	return &TerminalService{mgr: mgr}
}

// GetTerminalStatusById
func (s *TerminalService) GetTerminalStatus(ctx context.Context, id *v1.TerminalId (*pb.TerminalResponse, error) {
	terminal, err := s.mgr.GetTerminalByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	// check timeout
	isTimeout := time.Since(terminal.LastUpdated).Seconds() > float64(terminal.Timeout)
	status := terminal.Status
	if isTimeout {
		status = "offline"
	}

	return &v1.TerminalResponse{
		Id:      int32(terminal.ID),
		Status:  status,
		Timeout: int32(terminal.Timeout),
	}, nil
}

// UpdateTerminal
func (s *TerminalService) UpdateTerminal(ctx context.Context, req *pb.TerminalUpdateRequest) (*v1.TerminalResponse, error) {
	err := s.mgr.UpdateTerminal(ctx, int(req.Id), req.Status, int(req.Timeout))
	if err != nil {
		return nil, err
	}

	return &v1.TerminalResponse{
		Id:      req.Id,
		Status:  req.Status,
		Timeout: req.Timeout,
		Message: "Terminal updated successfully",
	}, nil
}

// SetTerminalTimeout
func (s *TerminalService) SetTerminalTimeout(ctx context.Context, req *v1.TerminalTimeoutRequest) (*pb.TerminalResponse, error) {
	err := s.mgr.SetTerminalTimeout(ctx, int(req.Id), int(req.Timeout))
	if err != nil {
		return nil, err
	}

	return &v1.TerminalResponse{
		Id:      req.Id,
		Timeout: req.Timeout,
		Status:  req.Status,
		Message: "Timeout updated successfully",
	}, nil
}
