package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DistributionKeeper interface {
	GetCommunityTax(sdk.Context) sdk.Dec
}

type distributionKeeper struct {
}

func (s *distributionKeeper) GetCommunityTax(sdk.Context) sdk.Dec {
	return sdk.ZeroDec()
}

type DistributionKeeperMock struct {
	s *distributionKeeper
}

func (s *DistributionKeeperMock) Mock() DistributionKeeper {
	return s.s
}

func NewDistributionKeeper() *DistributionKeeperMock {
	mock := &DistributionKeeperMock{s: &distributionKeeper{}}
	return mock
}
