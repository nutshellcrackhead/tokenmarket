package profile

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	protos "protos"
)

/**
PROFILE HANDLER MOCK
*/
type ProfileHandlerMock struct {
	mock.Mock
}

func (p *ProfileHandlerMock) Get(ctx context.Context, req *protos.ProfileGetRequest, rsp *protos.ProfileData) error {
	args := p.Called(ctx, req, rsp)
	return args.Error(0)
}
