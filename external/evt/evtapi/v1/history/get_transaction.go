package history

import (
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/client"
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/v1/chain"
)

type GetTransactionRequest struct {
	TransactionId string `json:"id"`
}

// type GetTransactionResult = chain.Transaction
type GetTransactionResult = chain.TransactionExtension

func (it *Instance) GetTransaction(transactionId string) (*GetTransactionResult, *client.ApiError) {
	response := &GetTransactionResult{}

	err := it.Client.Post(it.Path("get_transaction"), &GetTransactionRequest{transactionId}, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
