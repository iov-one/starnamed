package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type StakingKeeper interface {
	GetLastTotalPower(sdk.Context) sdk.Int
}

type stakingKeeper struct {
}

func (s *stakingKeeper) GetLastTotalPower(sdk.Context) sdk.Int {
	return sdk.OneInt()
}

type StakingKeeperMock struct {
	s *stakingKeeper
}

func (s *StakingKeeperMock) Mock() StakingKeeper {
	return s.s
}

func NewStakingKeeper() *StakingKeeperMock {
	mock := &StakingKeeperMock{s: &stakingKeeper{}}
	return mock
}
