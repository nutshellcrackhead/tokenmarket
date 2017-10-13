package handlers

import (
	"deposits/services"
	"golang.org/x/net/context"
	proto "protos"
)

type depositsHandler struct {
	proto.DepositsHandler
	storage services.DepositsStorageService
}

// Factory that creates new DepositsHandler instance
func NewDepositsHandler(storage services.DepositsStorageService) proto.DepositsHandler {
	return &depositsHandler{storage: storage}
}

// -----------MAIN MICROSERVICE METHODS------------

// DepositsHandler GetDeposits method implementation
func (handler *depositsHandler) GetDeposits(ctx context.Context, req *proto.GetDepositsRequest, rsp *proto.DepositsList) (err error) {
	userId := req.GetId()
	page := req.GetPage()
	state, err := handler.storage.GetDepositsList(userId, page)

	if err != nil {
		return err
	}

	rsp.Data = state

	return nil
}

func (handler *depositsHandler) CreateDeposit(ctx context.Context, req *proto.CreateDepositRequest, rsp *proto.DepositItem) error {
	userId := req.GetId()
	amount := req.GetAmount()
	currency := req.GetCurrency()

	err := handler.storage.CreateDeposit(rsp, userId, amount, currency)

	if err != nil {
		return err
	}

	return nil
}
