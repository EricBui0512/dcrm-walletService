package history

import (
	"fmt"

	"github.com/EricBui0512/dcrm-walletService/external/evt/evtapi/client"
	"github.com/EricBui0512/dcrm-walletService/external/evt/evtconfig"
)

type Instance struct {
	Client *client.Instance
	Config *evtconfig.Instance
}

func New(config *evtconfig.Instance, client *client.Instance) *Instance {
	return &Instance{
		Client: client,
		Config: config,
	}
}

func (it *Instance) Path(method string) string {
	return fmt.Sprintf("history/%v", method)
}
