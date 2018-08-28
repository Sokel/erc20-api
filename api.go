package api

import (
	"crypto/ecdsa"
	"math/big"

	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Tokener interface {
	Approve(key *ecdsa.PrivateKey, to string, amount *big.Int) (*types.Transaction, error)
	Transfer(key *ecdsa.PrivateKey, to string, amount *big.Int) (*types.Transaction, error)
	TransferFrom(key *ecdsa.PrivateKey, from string, to string, amount *big.Int) (*types.Transaction, error)
	BalanceOf(address string) (*big.Int, error)
	AllowanceOf(from string, to string) (*big.Int, error)
	TotalSupply() (*big.Int, error)
}

type API struct {
	client   *ethclient.Client
	gasPrice *big.Int

	tokenContract *StandardToken
}

func NewAPI(ethEndpoint string, gasPrice *big.Int, tokenContractAddress string) (*API, error) {
	client, err := ethclient.Dial(ethEndpoint)
	if err != nil {
		return nil, err
	}

	tokenContract, err := NewStandardToken(common.HexToAddress(tokenContractAddress), client)
	if err != nil {
		return nil, err
	}

	bch := &API{
		client:        client,
		gasPrice:      gasPrice,
		tokenContract: tokenContract,
	}
	return bch, nil
}

func (api *API) BalanceOf(ctx context.Context, address string) (*big.Int, error) {
	balance, err := api.tokenContract.BalanceOf(getCallOpts(ctx), common.HexToAddress(address))
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (api *API) AllowanceOf(ctx context.Context, from string, to string) (*big.Int, error) {
	allowance, err := api.tokenContract.Allowance(getCallOpts(ctx), common.HexToAddress(from), common.HexToAddress(to))
	if err != nil {
		return nil, err
	}
	return allowance, nil
}

func (api *API) Approve(ctx context.Context, key *ecdsa.PrivateKey, to string, amount *big.Int) (*types.Transaction, error) {
	opts := getTxOpts(key, 48000)
	tx, err := api.tokenContract.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return nil, err
	}
	return tx, err
}

func (api *API) Transfer(ctx context.Context, key *ecdsa.PrivateKey, to string, amount *big.Int) (*types.Transaction, error) {
	opts := getTxOpts(key, 40000)
	tx, err := api.tokenContract.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return nil, err
	}
	return tx, err
}

func (api *API) TransferFrom(ctx context.Context, key *ecdsa.PrivateKey, from string, to string, amount *big.Int) (*types.Transaction, error) {
	opts := getTxOpts(key, 50000)
	tx, err := api.tokenContract.TransferFrom(opts, common.HexToAddress(from), common.HexToAddress(to), amount)
	if err != nil {
		return nil, err
	}
	return tx, err
}

func (api *API) TotalSupply(ctx context.Context) (*big.Int, error) {
	totalSupply, err := api.tokenContract.TotalSupply(getCallOpts(ctx))
	if err != nil {
		return nil, err
	}
	return totalSupply, nil
}

func getTxOpts(key *ecdsa.PrivateKey, gasLimit uint64) *bind.TransactOpts {
	opts := bind.NewKeyedTransactor(key)
	opts.GasLimit = gasLimit
	return opts
}

func getCallOpts(ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{
		Context: ctx,
		Pending: false,
	}
}
