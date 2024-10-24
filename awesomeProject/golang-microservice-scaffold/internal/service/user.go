package service

import (
	"context"
	v1 "example/api/user/v1"
	"example/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserService is the service interface for other services or users to call
//
// You should embed the UnimplementedXXServer struct into your service struct and override the RPC methods,
// so that the servers would automatically provide the interfaces
type UserService struct {
	v1.UnimplementedUserManagementServer
	// mgr is the business layer operation collection, which implements the service interface
	mgr *biz.UserManager
}

func NewUserService(mgr *biz.UserManager) *UserService {
	return &UserService{mgr: mgr}
}

func (s *UserService) AddUser(ctx context.Context, usr *v1.User) (empty *emptypb.Empty, err error) {
	if valid := usr.Validate(); valid != nil {
		return nil, v1.ErrorMalformedInput("Malformed user information: %v", valid)
	}
	err = s.mgr.Add(ctx, usr)
	return
}
func (s *UserService) UpdateUser(ctx context.Context, usr *v1.User) (empty *emptypb.Empty, err error) {
	if valid := usr.Validate(); valid != nil {
		return nil, v1.ErrorMalformedInput("Malformed user information: %v", valid)
	}
	err = s.mgr.Update(ctx, usr)
	return
}
func (s *UserService) FindUserByName(ctx context.Context, name *v1.UserName) (usr *v1.User, err error) {
	if valid := name.Validate(); valid != nil {
		return nil, v1.ErrorMalformedInput("Malformed name: %v", valid)
	}
	return s.mgr.GetByName(ctx, name.Name)
}
func (s *UserService) RemoveUserById(ctx context.Context, uid *v1.UserId) (empty *emptypb.Empty, err error) {
	err = s.mgr.RemoveById(ctx, uid.Id)
	return
}
