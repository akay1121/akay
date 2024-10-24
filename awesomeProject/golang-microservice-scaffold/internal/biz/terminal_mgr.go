package biz

import (
	"context"
	v1 "example/api/terminal"

	"example/internal/constant"
	"time"
)

// What we should do here includes:
//   - Create a new entity struct type that the repository operates on
//   - Define the repository interface
//   - Implement the business logic
// You shall notice that the code below only focuses on the business logic rather than cope with the interaction
// with the exposed service APIs. You should never deal with HTTP requests or something likewise here.

// User is just an alias for the entity type.
//
// Here we just use the type directly for convenience, but we strongly recommend you write your own entity
// structure instead, since entities in the domain layer are slightly different from ones in the data access layer,
// or infrastructure layer in a typical DDD repository.
type Terminal = v1.Terminal

type TerminalRepository interface {
	GetTerminalByID(ctx context.Context, id int) (*Terminal, error)
	GetTerminalStatusByID(ctx context.Context, id int) (string, error)
	UpdateTerminal(ctx context.Context, terminal Terminal) error
	SetTerminalTimeout(ctx context.Context, id int, timeout int) error
	IsTerminalExist(ctx context.Context, id int) (bool, error)
}

// UserManager is where the business logic resides. It encapsulates the repository inside and provides intuitive
// operations to help the upper layers only concentrate on the business logic instead of manipulating the repository.
type TerminalManager struct {
	repo    TerminalRepository
	timeout time.Duration
}

func NewTerminalManager(repo TerminalRepository) *TerminalManager {
	return &TerminalManager{repo: repo}
}

func (m *TerminalManager) Update(ctx context.Context, terminal *Terminal) error {
	return m.repo.UpdateTerminal(ctx, terminal)
}

func (m *TerminalManager) GetTerminalById(ctx context.Context, id int) (terminal *Terminal, err error) {
	var ext bool
	if ext, err = m.repo.IsTerminalExist(ctx, id); err != nil {
		return
	}
	if !ext {
		return nil, v1.ErrorTerminalNotFound("There is no such Terminal id %v", id)
	}
	return m.repo.GetTerminalByID(ctx, id)
}
func (m *TerminalManager) GetTerminalStatus(ctx context.Context, id int) (string, error) {
	terminal, err := m.repo.GetTerminalByID(ctx, id)
	if err != nil {
		return "", err
	}
	// if timeout status become offline
	if time.Since(terminal.LastUpdated) > m.timeout {
		return constant.TerminalStatusOffline, nil
	}
	return terminal.Status, nil
}
func (m *TerminalManager) SetTimeOut(ctx context.Context, id int, timeout int) error {
	return m.repo.SetTerminalTimeout(ctx, id, timeout)
}
