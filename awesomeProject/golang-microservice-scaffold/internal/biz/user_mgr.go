package biz

import (
	"context"
	v1 "example/api/user/v1"
	"example/internal/ent"
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
type User = v1.User

// UserRepository represents the interface of operating the entities stored in the database, no matter where it is,
// local disk or remote storage.
//
// The implementation of UserRepository should be placed in the [example/internal/data] package, in order not to
// interfere with the business of this package. It just serves as the entry to the real implementation.
//
// You may confuse the repository with the business logic - they both operate on the entity type. Their main ambition
// is, however, totally different. Repositories do not care about the validity of the passed parameters, it just
// does CRUD operations; business logic (or service), in contrast, would think over the design, and do more operations.
// For example, a service may increment a user's points if commodities are purchased successfully - it operates the
// user and commodity repositories at the same time.
type UserRepository interface {
	Add(ctx context.Context, user *User) error
	Remove(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	FindByName(ctx context.Context, name string) (*User, error)
	FindById(ctx context.Context, id int64) (*User, error)
	FindChildrenByParentId(ctx context.Context, id int64) ([]*User, error)
	RecoverById(ctx context.Context, id int64) error
	IsNameExist(ctx context.Context, name string) (bool, error)
}

// UserManager is where the business logic resides. It encapsulates the repository inside and provides intuitive
// operations to help the upper layers only concentrate on the business logic instead of manipulating the repository.
type UserManager struct {
	repo UserRepository
}

func NewUserManager(repo UserRepository) *UserManager {
	return &UserManager{repo: repo}
}

func (m *UserManager) Add(ctx context.Context, user *User) error {
	return m.repo.Add(ctx, user)
}

func (m *UserManager) RemoveById(ctx context.Context, id int64) (err error) {
	var usr *User
	if usr, err = m.repo.FindById(ctx, id); err != nil {
		switch {
		case ent.IsNotFound(err):
			return v1.ErrorUserNotFound("Cannot find the specified user with id %v", id)
		}
		return
	}
	return m.repo.Remove(ctx, usr)
}

func (m *UserManager) Update(ctx context.Context, user *User) (err error) {
	return m.repo.Update(ctx, user)
}

func (m *UserManager) GetByName(ctx context.Context, name string) (usr *User, err error) {
	var ext bool
	if ext, err = m.repo.IsNameExist(ctx, name); err != nil {
		return
	}
	if !ext {
		return nil, v1.ErrorUserNotFound("There is no such user named %v", name)
	}
	return m.repo.FindByName(ctx, name)
}
