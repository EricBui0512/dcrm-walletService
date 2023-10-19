package v1

import (
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/client"
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/v1/chain"
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/v1/evt"
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/v1/evt_link"
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/v1/history"
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtconfig"
)

type Instance struct {
	Chain   *chain.Instance
	EvtLink *evt_link.Instance
	Evt     *evt.Instance
	History *history.Instance
}

func New(config *evtconfig.Instance, client *client.Instance) *Instance {
	return &Instance{
		Chain:   chain.New(config, client),
		Evt:     evt.New(config, client),
		History: history.New(config, client),
	}
}
