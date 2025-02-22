package evt

import (
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/client"
)

type GetTokensRequest struct {
	Domain string `json:"domain"`
	Skip   int    `json:"skip"`
	Take   int    `json:"take"`
}

type GetTokensResult = []interface{}

func (it *Instance) GetTokens(domainName string, skip int, take int) (*GetTokensResult, *client.ApiError) {
	response := &GetTokensResult{}

	err := it.client.Post(it.path("get_tokens"), &GetTokensRequest{domainName, skip, take}, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
