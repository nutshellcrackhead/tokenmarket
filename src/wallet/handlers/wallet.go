package handlers

import (
	"golang.org/x/net/context"
	proto "protos"
	"wallet/services"
)

type walletHandler struct {
	proto.WalletHandler
	storage services.WalletStorageService
}

// Factory that creates new WalletHandler instance
func NewWalletHandler(storage services.WalletStorageService) proto.WalletHandler {
	return &walletHandler{storage: storage}
}

// -----------MAIN MICROSERVICE METHODS------------

// WalletHandler Login method implementation
func (handler *walletHandler) Status(ctx context.Context, req *proto.WalletStatusRequest, rsp *proto.WalletStatusResponse) (err error) {
	userId := req.GetId()
	state, err := handler.storage.GetWalletStatus(userId)

	if err != nil {
		return err
	}

	rsp.State = state

	return nil
}

func (handler *walletHandler) OperationsList(ctx context.Context, req *proto.WalletOperationsListRequest, rsp *proto.WalletOperationsList) error {
	userId := req.GetId()
	page := req.GetPage()

	operations, err := handler.storage.GetOperationsList(userId, page)

	if err != nil {
		return err
	}

	rsp.Data = operations

	return nil
}
