package handlers

import (
	"golang.org/x/net/context"
	proto "protos"
	"referral/services"
)

type referralHandler struct {
	proto.ReferralHandler
	storage services.ReferralStorageService
}

// Factory that creates new ReferralHandler instance
func NewReferralHandler(storage services.ReferralStorageService) proto.ReferralHandler {
	return &referralHandler{storage: storage}
}

// -----------MAIN MICROSERVICE METHODS------------

// Referral Status method implementation
func (handler *referralHandler) List(ctx context.Context, req *proto.ReferralListRequest, rsp *proto.ReferralPartnerList) (err error) {
	userId := req.GetId()
	page := req.GetPage()

	data := []*proto.ReferralPartner{}

	err = handler.storage.GetReferralStatus(&data, userId, page)

	if err != nil {
		return err
	}

	rsp.Data = data

	return nil
}
