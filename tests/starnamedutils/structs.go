package starnamedutils

import (
	"context"
	"fmt"
	"time"

	tools "github.com/iov-one/starnamed/tests/starnametesttools"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

// Module: Starname

const (
	typeDomain       string = "domain"
	typeAccount      string = "account"
	domainKindOpen   string = "open"
	domainKindClosed string = "closed"
)

type StarnameDomain struct {
	chain          *cosmos.CosmosChain
	name           string
	is_open_domain bool
	owner          ibc.Wallet
	command        *tools.Command
	ctx            context.Context
}

func (d *StarnameDomain) Command(command *tools.Command, ctx context.Context) *StarnameDomain {
	return &StarnameDomain{
		chain:          d.chain,
		name:           d.name,
		is_open_domain: d.is_open_domain,
		owner:          d.owner,
		command:        command,
		ctx:            ctx,
	}
}

func (d *StarnameDomain) GetName() string {
	return d.name
}

func (d *StarnameDomain) GetOwner() ibc.Wallet {
	return d.owner
}

func (d *StarnameDomain) DomainKind() string {
	if d.is_open_domain {
		return domainKindOpen
	}
	return domainKindClosed
}

func (d *StarnameDomain) Delete() error {
	std_out, std_err, err := d.command.Tx(d.owner, true, false).SetArgs("starname", "domain-delete", "--domain", d.name).Exec(d.ctx)

	if err != nil {
		return err
	}

	if string(std_err) != "" {
		return fmt.Errorf("error: %s", std_err)
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		return fmt.Errorf("no output")
	}

	return nil
}

func (d *StarnameDomain) TransferOwnership(new_owner ibc.Wallet) error {
	std_out, std_err, err := d.command.Tx(d.owner, true, false).SetArgs("starname", "domain-transfer", "--domain", d.name, "--new-owner", new_owner.FormattedAddress()).Exec(d.ctx)

	if err != nil {
		return err
	}

	if string(std_err) != "" {
		return fmt.Errorf("error: %s", std_err)
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		return fmt.Errorf("no output")
	}

	d.owner = new_owner

	return nil

}

func (d *StarnameDomain) Renew() error {
	std_out, std_err, err := d.command.Tx(d.owner, true, false).SetArgs("starname", "domain-renew", "--domain", d.name).Exec(d.ctx)

	if err != nil {
		return err
	}

	if string(std_err) != "" {
		return fmt.Errorf("error: %s", std_err)
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		return fmt.Errorf("no output")
	}

	return nil
}

func (d *StarnameDomain) Type() string {
	return typeDomain
}

func (d *StarnameDomain) Escrow(price int64, denom string) (escrow Escrow, err error) {

	price_str := fmt.Sprintf("%d%s", price, denom)

	//Time must be in RFC3339 format, will be 1 hour from now
	time := time.Now().Add(time.Hour).Format(time.RFC3339)

	std_out, std_err, err := d.command.Tx(d.owner, true, false).SetArgs("starname", "domain-escrow-create", "--domain", d.name, "--price", price_str, "--expiration", time).Exec(d.ctx)
	escrow = Escrow{}

	if err != nil {
		return
	}

	if string(std_err) != "" {
		err = fmt.Errorf("error: %s", std_err)
		return
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		err = fmt.Errorf("no output")
		return
	}

	return Escrow{
		chain:  d.chain,
		owner:  d.owner,
		object: d,
		name:   d.name,
		price:  price,
		denom:  denom,
	}, nil
}

type StarnameAccount struct {
	chain   *cosmos.CosmosChain
	name    string
	domain  StarnameDomain
	owner   ibc.Wallet
	command *tools.Command
	ctx     context.Context
}

func (a *StarnameAccount) Command(command *tools.Command, ctx context.Context) *StarnameAccount {
	return &StarnameAccount{
		chain:   a.chain,
		name:    a.name,
		domain:  a.domain,
		owner:   a.owner,
		command: command,
		ctx:     ctx,
	}
}

func (a *StarnameAccount) GetName() string {
	return a.name
}

func (a *StarnameAccount) GetOwner() ibc.Wallet {
	return a.owner
}

func (a *StarnameAccount) GetDomain() StarnameDomain {
	return a.domain
}

func (a *StarnameAccount) Delete() error {

	std_out, std_err, err := a.command.Tx(a.owner, true, false).SetArgs("starname", "account-delete", "--domain", a.domain.GetName(), "--name", a.name).Exec(a.ctx)

	if err != nil {
		return err
	}

	if string(std_err) != "" {
		return fmt.Errorf("error: %s", std_err)
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		return fmt.Errorf("no output")
	}

	return nil
}

func (a *StarnameAccount) TransferOwnership(new_owner ibc.Wallet) error {

	std_out, std_err, err := a.command.Tx(a.owner, true, false).SetArgs("starname", "account-transfer", "--domain", a.domain.GetName(), "--name", a.name, "--new-owner", new_owner.FormattedAddress()).Exec(a.ctx)

	if err != nil {
		return err
	}

	if string(std_err) != "" {
		return fmt.Errorf("error: %s", std_err)
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		return fmt.Errorf("no output")
	}

	a.owner = new_owner

	return nil
}

func (d *StarnameAccount) Renew() error {

	std_out, std_err, err := d.command.Tx(d.owner, true, false).SetArgs("starname", "account-renew", "--domain", d.domain.GetName(), "--name", d.name).Exec(d.ctx)

	if err != nil {
		return err
	}

	if string(std_err) != "" {
		return fmt.Errorf("error: %s", std_err)
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		return fmt.Errorf("no output")
	}

	return nil
}

func (a *StarnameAccount) Type() string {
	return typeAccount
}

func (a *StarnameAccount) Escrow(price int64, denom string) (escrow Escrow, err error) {

	price_str := fmt.Sprintf("%d%s", price, denom)

	//Time must be in RFC3339 format, will be 1 hour from now
	time := time.Now().Add(time.Hour).Format(time.RFC3339)

	std_out, std_err, err := a.command.Tx(a.owner, true, false).SetArgs("starname", "account-escrow-create", "--domain", a.domain.GetName(), "--name", a.name, "--price", price_str, "--expiration", time).Exec(a.ctx)

	if err != nil {
		return
	}

	if string(std_err) != "" {
		err = fmt.Errorf("error: %s", std_err)
		return
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		err = fmt.Errorf("no output")
		return
	}

	return Escrow{
		chain:  a.chain,
		owner:  a.owner,
		object: a,
		name:   a.name,
		price:  price,
		denom:  denom,
	}, nil
}

// Module: Escrow
type Escrow struct {
	chain   *cosmos.CosmosChain
	command *tools.Command
	ctx     context.Context
	owner   ibc.Wallet
	object  Escrowobject
	name    string
	price   int64
	denom   string
	id      string
}

type Tradable interface {
	Escrow(price int64, denom string) (Escrow, error)
	Type() string
}

type Escrowobject interface {
	Type() string
}

func (e *Escrow) Command(command *tools.Command) *Escrow {
	return &Escrow{
		chain:   e.chain,
		command: command,
		owner:   e.owner,
		object:  e.object,
		name:    e.name,
		price:   e.price,
		denom:   e.denom,
		id:      e.id,
	}
}

func (e *Escrow) Delete() (err error) {
	// This the cli command: refound

	std_out, std_err, err := e.command.Tx(e.owner, true, false).SetArgs("escrow", "refund", e.id).Exec(e.ctx)

	if err != nil {
		return
	}

	if string(std_err) != "" {
		err = fmt.Errorf("error: %s", std_err)
		return
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		err = fmt.Errorf("no output")
		return
	}

	return nil
}

func (e *Escrow) Buy(new_owner ibc.Wallet) (err error) {
	// This the cli command: Transfer
	price_str := fmt.Sprintf("%d%s", e.price, e.denom)
	std_out, std_err, err := e.command.Tx(new_owner, true, false).SetArgs("escrow", "transfer", e.id, price_str).Exec(e.ctx)

	if err != nil {
		return
	}

	if string(std_err) != "" {
		err = fmt.Errorf("error: %s", std_err)
		return
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		err = fmt.Errorf("no output")
		return
	}

	e.owner = new_owner

	return
}

func (e *Escrow) UpdatePrice(new_price int64, denom string) (err error) {
	// This the cli command: Update

	price_str := fmt.Sprintf("%d%s", e.price, e.denom)
	std_out, std_err, err := e.command.Tx(e.owner, true, false).SetArgs("escrow", "update", e.id, "--price", price_str).Exec(e.ctx)

	if err != nil {
		return
	}

	if string(std_err) != "" {
		err = fmt.Errorf("error: %s", std_err)
		return
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		err = fmt.Errorf("no output")
		return
	}

	e.price = new_price
	e.denom = denom
	return
}

// Factory

func NewStarnameDomain(cli *tools.Command, ctx context.Context, chain *cosmos.CosmosChain, name string, owner ibc.Wallet, is_open_domain bool) (*StarnameDomain, error) {

	domainKind := domainKindClosed
	if is_open_domain {
		domainKind = domainKindOpen
	}

	domainName := name

	if domainName == "" {
		domainName = randomString(10)
	}

	std_out, std_err, err := cli.Tx(owner, true, false).SetArgs("starname", "domain-register", "--domain", domainName, "--type", domainKind).Exec(ctx)

	if err != nil {
		return &StarnameDomain{}, err
	}

	if string(std_err) != "" {
		return &StarnameDomain{}, fmt.Errorf("error: %s", std_err)
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator
	if string(std_out) == "" {
		return &StarnameDomain{}, fmt.Errorf("no output")
	}

	json, err := StringUnmarshal(string(std_out))

	if err != nil || json == nil {
		return &StarnameDomain{}, err
	}

	domain := &StarnameDomain{
		chain:          chain,
		name:           domainName,
		is_open_domain: is_open_domain,
		owner:          owner,
	}

	return domain, nil
}

func NewStarnameAccount(cli *tools.Command, ctx context.Context, chain *cosmos.CosmosChain, name string, owner ibc.Wallet, domain StarnameDomain) (account *StarnameAccount, err error) {

	// if domain.DomainKind() == "closed" && domain.GetOwner().FormattedAddress() != owner.FormattedAddress() {
	// 	return StarnameAccount{}, fmt.Errorf("only the owner of a closed domain can create an account")
	// }

	accountName := name

	if accountName == "" {
		accountName = randomString(10)
	}

	std_out, std_err, err := cli.Tx(owner, true, false).SetArgs("starname", "account-register", "--domain", domain.GetName(), "--name", accountName).Exec(context.Background())

	if err != nil {
		return
	}

	if string(std_err) != "" {
		err = fmt.Errorf("error: %s", std_err)
		return
	}

	// turn the std_out into a json object
	// TODO: Implement a json validator

	if string(std_out) == "" {
		err = fmt.Errorf("no output")
		return
	}

	return &StarnameAccount{
		chain:  chain,
		name:   accountName,
		domain: domain,
		owner:  domain.owner,
	}, nil
}
