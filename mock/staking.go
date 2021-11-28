package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type StakingKeeper interface {
	GetLastTotalPower(sdk.Context) sdk.Int
	TokensFromConsensusPower(sdk.Context, int64) sdk.Int
}

type stakingKeeper struct {
}

func (s *stakingKeeper) GetLastTotalPower(sdk.Context) sdk.Int {
	return sdk.OneInt()
}

func (s *stakingKeeper) TokensFromConsensusPower(_ sdk.Context, power int64) sdk.Int {
	return sdk.TokensFromConsensusPower(power, sdk.NewInt(int64(1e6)))
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
