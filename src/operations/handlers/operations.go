package handlers

import (
	"errors"
	"golang.org/x/net/context"
	"helpers/values"
	"operations/services"
	proto "protos"
)

type operationsHandler struct {
	proto.OperationsHandler
	storage services.OperationsStorageService
}

// Factory that creates new OperationsHandler instance
func NewOperationsHandler(storage services.OperationsStorageService) proto.OperationsHandler {
	return &operationsHandler{storage: storage}
}

// -----------MAIN MICROSERVICE METHODS------------

// Operations method implementation
func (handler *operationsHandler) CreatePayment(ctx context.Context, req *proto.CreatePaymentRequest, rsp *proto.CreatePaymentResponse) (err error) {
	id := req.GetId()
	method := req.GetMethod()
	amount := req.GetAmount()
	currencyName := req.GetCurrency()
	merchantAccount := values.GetMerchant(method)

	if merchantAccount == "" {
		err = errors.New("Unknown merchant")
		return
	}

	err = handler.storage.CreatePayment(rsp, id, method, amount, currencyName, "token_topup")

	if err != nil {
		return
	}

	rsp.MerchantAccount = &merchantAccount

	return nil
}

func (handler *operationsHandler) CreatePayout(ctx context.Context, req *proto.CreatePayoutRequest, rsp *proto.CreatePayoutResponse) (err error) {
	id := req.GetId()
	method := req.GetMethod()
	amount := req.GetAmount()
	currencyName := req.GetCurrency()
	merchantAccount := values.GetMerchant(method)
	account := req.GetAccount()

	if merchantAccount == "" {
		err = errors.New("Unknown merchant")
		return
	}

	err = handler.storage.CreatePayout(id, method, amount, currencyName, account, "token_payout_operation")

	if err != nil {
		return
	}

	okStatus := "ok"

	rsp.Status = &okStatus

	return nil
}

func (handler *operationsHandler) ActivateUser(ctx context.Context, req *proto.CreatePaymentRequest, rsp *proto.CreatePaymentResponse) (err error) {
	id := req.GetId()
	method := req.GetMethod()
	amount := req.GetAmount()
	currencyName := req.GetCurrency()
	merchantAccount := values.GetMerchant(method)

	if merchantAccount == "" {
		err = errors.New("Unknown merchant")
		return
	}

	err = handler.storage.CreatePayment(rsp, id, method, amount, currencyName, "token_account_activation")

	if err != nil {
		return
	}

	err = handler.storage.CreateAccountActivation(rsp.GetId())

	if err != nil {
		return
	}

	rsp.MerchantAccount = &merchantAccount

	return nil
}

func (handler *operationsHandler) UpdatePaymentStatus(ctx context.Context, req *proto.UpdatePaymentRequest, rsp *proto.UpdatePaymentResponse) (err error) {
	paymentId := req.GetOperationId()
	token := req.GetToken()
	amount := req.GetAmount()
	currency := req.GetCurrency()
	method := req.GetMethod()
	status := req.GetStatus()

	operationId, err := handler.storage.ValidatePayment(paymentId, token, amount, currency, method)

	if err != nil {
		return
	}

	if operationId == 0 {
		err = errors.New("Not valid payment")
		return
	}

	err = handler.storage.UpdatePayment(operationId, status)

	if err != nil {
		return
	}

	go handler.storage.UnlockAccount(operationId)
	go handler.storage.TopUpWallet(operationId, status)

	rsp.Status = &status

	return nil
}
