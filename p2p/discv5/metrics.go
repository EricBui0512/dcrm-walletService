package discv5

import "github.com/EricBui0512/dcrm-walletService/p2p/metrics"

var (
	ingressTrafficMeter = metrics.NewRegisteredMeter("discv5/InboundTraffic", nil)
	egressTrafficMeter  = metrics.NewRegisteredMeter("discv5/OutboundTraffic", nil)
)
