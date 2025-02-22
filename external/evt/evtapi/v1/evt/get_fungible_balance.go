package evt

import (
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/client"
)

type GetFungibleBalanceRequest struct {
	SymId   uint   `json:"sym_id"`
	Address string `json:"address"`
}

type GetFungibleBalanceResult = []string

func (it *Instance) GetFungibleBalance(symid uint, address string) (*GetFungibleBalanceResult, *client.ApiError) {
	response := &GetFungibleBalanceResult{}

	err := it.client.Post(it.path("get_fungible_balance"), &GetFungibleBalanceRequest{
		SymId:   symid,
		Address: address,
	}, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
