package data

import (
	"context"
	v1 "example/api/user/v1"
	"example/internal/biz"
	"example/internal/ent"
	"example/internal/ent/user"
	"github.com/jinzhu/copier"
)

// userRepo implements the interface [biz.UserRepository] described in the package [example/internal/biz].
//
// You should always name the implementation struct with a lowercase letter to avoid exposing the implementation
// to other packages.
type userRepo struct {
	db    *Data
	cache *Cache
}

// NewUserRepository creates a new user repository implementation instance and returns the interface value.
//
// The reason for returning the interface instead of the instance itself is to ensure SOLID principles -
// you should always depend on the interface rather than the underlying implementation.
func NewUserRepository(database *Data, cache *Cache) biz.UserRepository {
	return &userRepo{db: database, cache: cache}
}

// If you are only interested in the layout of the example project, you can skip the following implementation details.
// Otherwise, you should read the comments to understand the internal design of the example project.

var (
	// Redis key used to store the usage of usernames
	keyUsername = "user:names"
)

func convertToBizUser(u *ent.User) (usr *biz.User, err error) {
	usr = &biz.User{}
	if err = copier.Copy(usr, u); err != nil {
		return
	}
	// Make sure the input only fields are not included
	usr.Password = nil
	usr.Secret = nil
	return
}

func (r *userRepo) Add(ctx context.Context, u *biz.User) (err error) {
	// The topmost user has a parent id of -1; any other user should have a valid parent id that exists in the DB
	if u.ParentId != -1 {
		var parentExists bool
		// Find if there exists a user with the specified id
		if parentExists, err = r.db.Client.User.Query().Where(user.IDEQ(u.ParentId)).Exist(ctx); err != nil {
			return
		}
		if !parentExists {
			// Here we return the generated error
			return v1.ErrorUserNotFound("Invalid parent id specified")
		}
	}
	_, err = r.db.Client.User.Create().
		SetParentID(u.ParentId).
		SetName(u.Name).
		SetNickname(u.Nickname).
		SetPassword(*u.Password).
		SetEmail(u.Email).
		Save(ctx)
	if err == nil {
		r.cache.Client.BFAdd(ctx, keyUsername, u.Name)
	}
	return
}
func (r *userRepo) Remove(ctx context.Context, u *biz.User) (err error) {
	var usr *ent.User
	if usr, err = r.db.Client.User.Query().Where(user.IDEQ(u.Id)).First(ctx); err != nil {
		return
	}
	// We would do nothing but set the deleted flag
	return r.db.Client.User.Update().Where(user.IDEQ(usr.ID)).SetDeleted(true).Exec(ctx)
}
func (r *userRepo) Update(ctx context.Context, u *biz.User) error {
	return r.db.Client.User.Update().Where(user.IDEQ(u.Id)).Exec(ctx)
}
func (r *userRepo) FindByName(ctx context.Context, name string) (usr *biz.User, err error) {
	var u *ent.User
	if u, err = r.db.Client.User.Query().Where(user.NameEQ(name)).First(ctx); err != nil {
		return
	}
	return convertToBizUser(u)
}
func (r *userRepo) FindById(ctx context.Context, id int64) (usr *biz.User, err error) {
	var u *ent.User
	if u, err = r.db.Client.User.Query().Where(user.IDEQ(id)).First(ctx); err != nil {
		return
	}
	return convertToBizUser(u)
}
func (r *userRepo) FindChildrenByParentId(ctx context.Context, id int64) (children []*biz.User, err error) {
	children = make([]*biz.User, 0, 3)
	var u *ent.User
	if u, err = r.db.Client.User.Query().Where(user.IDEQ(id)).First(ctx); err != nil {
		return
	}
	var rawChildren []*ent.User
	if rawChildren, err = r.db.Client.User.QueryChildren(u).All(ctx); err != nil {
		return
	}
	var bizChild *biz.User
	for _, child := range rawChildren {
		if bizChild, err = convertToBizUser(child); err != nil {
			return
		}
		children = append(children, bizChild)
	}
	return
}
func (r *userRepo) RecoverById(ctx context.Context, id int64) (err error) {
	var usr *ent.User
	if usr, err = r.db.Client.User.Query().Where(user.IDEQ(id)).First(ctx); err != nil {
		return
	}
	return r.db.Client.User.Update().Where(user.IDEQ(usr.ID)).SetDeleted(false).Exec(ctx)
}
func (r *userRepo) IsNameExist(ctx context.Context, name string) (bool, error) {
	// Check if name already exists using Bloom filter
	if r.cache.Client.BFExists(ctx, keyUsername, name).Val() { // Might contain the name already
		ex, err := r.db.Client.User.Query().Where(user.NameEQ(name)).Exist(ctx) // Check the database for sure
		return !ex, err                                                         // If the name already exists, return false
	} else { // Bloom filter indicates the value does not appear, it must be non-existent
		return true, nil
	}
}
