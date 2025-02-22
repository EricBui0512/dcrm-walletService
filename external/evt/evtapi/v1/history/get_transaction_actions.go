package history

import "github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/client"

type GetTransactionActionsRequest struct {
	TransactionId string `json:"id"`
}

type GetTransactionIdsForBlockResult = []string

func (it *Instance) GetTransactionActions(transactionId string) (*GetTransactionActionsRequest, *client.ApiError) {
	response := &GetTransactionActionsRequest{}

	err := it.Client.Post(it.Path("get_transaction_actions"), &GetTransactionActionsRequest{transactionId}, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
