package service

import (
	"ShunFengParcel/internal/biz"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGreeterService,
	NewOrderServiceImpl,
	biz.NewPricingUsecase,
)
