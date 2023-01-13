package ports

import (
	"context"

	"github.com/BON4/gofeed/internal/accounts/app"
	"github.com/BON4/gofeed/internal/common/errors"
	"github.com/BON4/gofeed/internal/common/genproto/accounts"
	"github.com/golang/protobuf/ptypes/empty"
)

type GrpcServer struct {
	app *app.Application
}

func NewGrpcServer(a *app.Application) *GrpcServer {
	return &GrpcServer{
		app: a,
	}
}

func (g *GrpcServer) ChangePassword(context.Context, *accounts.ChangePasswordRequest) (*empty.Empty, error) {
	return &empty.Empty{}, errors.NewNotImplementedError("not implemented", "not-implemented")
}

func (g *GrpcServer) ActivateAccount(context.Context, *accounts.ActivateAccountRequest) (*empty.Empty, error) {
	return &empty.Empty{}, errors.NewNotImplementedError("not implemented", "not-implemented")
}
