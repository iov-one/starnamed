package starnametesttools

import (
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

type StarnameIBCWallet struct {
	Wallet       ibc.Wallet
	Chain        *cosmos.CosmosChain
	StarDomains  []string
	StarAccounts []string
}

// ############################## Wallet ##############################
func (w *StarnameIBCWallet) KeyName() string {
	return w.Wallet.KeyName()
}

func (w *StarnameIBCWallet) FormattedAddress() string {
	return w.Wallet.FormattedAddress()
}

func (w *StarnameIBCWallet) Mnemonic() string {
	return w.Wallet.Mnemonic()
}

func (w *StarnameIBCWallet) Address() []byte {
	return w.Wallet.Address()
}

func (w *StarnameIBCWallet) GetIBCWallet() ibc.Wallet {
	return w.Wallet
}
