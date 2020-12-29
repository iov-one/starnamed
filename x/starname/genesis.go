package starname

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/starname/types"
)

// NewGenesisState builds a genesis state including the domains provided
func NewGenesisState(domains []types.Domain, accounts []types.Account) types.GenesisState {
	return types.GenesisState{Domains: domains, Accounts: accounts}
}

// ValidateGenesis validates a genesis state
// checking for domain validity and no domain name repetitions
func ValidateGenesis(data types.GenesisState) error {
	namesSet := make(map[string]struct{}, len(data.Domains))
	for _, domain := range data.Domains {
		if _, ok := namesSet[domain.Name]; ok {
			return fmt.Errorf("domain name %s declared twice", domain.Name)
		}
		namesSet[domain.Name] = struct{}{}
		if err := validateDomain(domain); err != nil {
			return err
		}
	}
	return nil
}

// DefaultGenesisState creates an empty genesis state for the domain module
func DefaultGenesisState() types.GenesisState {
	return types.GenesisState{Domains: []types.Domain{}, Accounts: []types.Account{}}
}

// InitGenesis builds a state from GenesisState
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	// insert domains
	ds := keeper.DomainStore(ctx)
	for _, domain := range data.Domains {
		// create domains
		ds.Create(&domain)
	}
	// insert accounts
	as := keeper.AccountStore(ctx)
	for _, account := range data.Accounts {
		as.Create(&account)
	}
}

// ExportGenesis saves the state of the domain module
func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	// domains
	cursor, err := k.DomainStore(ctx).Query().Where().Index(0x2).Equals([]byte{0x2}).Do() // HARD-CODE in conjunction with keeper.go
	if err != nil {
		panic(err)
	}
	var domains []types.Domain
	domain := new(types.Domain)
	for ; cursor.Valid(); cursor.Next() {
		err = cursor.Read(domain)
		if err != nil {
			panic(err)
		}
		domains = append(domains, *domain)
	}

	// accounts
	cursor, err = k.AccountStore(ctx).Query().Where().Index(0x1).Equals([]byte{0x1}).Do() // HARD-CODE in conjunction with keeper.go
	if err != nil {
		panic(err)
	}
	var accounts []types.Account
	account := new(types.Account)
	for ; cursor.Valid(); cursor.Next() {
		err = cursor.Read(account)
		if err != nil {
			panic(err)
		}
		accounts = append(accounts, *account)
	}

	return &types.GenesisState{
		Domains:  domains,
		Accounts: accounts,
	}
}

// validateDomain checks if a domain is valid or not
func validateDomain(d types.Domain) error {
	// TODO fill
	return nil
}
