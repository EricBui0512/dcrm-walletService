// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// bootnode runs a bootstrap node for the Ethereum Discovery Protocol.
package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/BurntSushi/toml"
	//"github.com/fusion/go-fusion/cmd/utils"
	"github.com/EricBui0512/dcrm-walletService/crypto"
	"github.com/EricBui0512/dcrm-walletService/p2p/discover"
	"github.com/EricBui0512/dcrm-walletService/p2p/discv5"
	"github.com/EricBui0512/dcrm-walletService/p2p/nat"
	"github.com/EricBui0512/dcrm-walletService/p2p/netutil"
)

func main() {
	var (
		groupNum      = flag.Uint("group", uint(0), "group Number") //0:sdk, 1:one group, 2:two groups dcrm xp
		groupNodesNum = flag.Uint("nodes", uint(0), "nodes Number in some group, must > 0")
		listenAddr    = flag.String("addr", "", "listen address")
		genKey        = flag.String("genkey", "", "generate a node key")
		writeAddr     = flag.Bool("writeaddress", false, "write out the node's pubkey hash and quit")
		nodeKeyFile   = flag.String("nodekey", "", "private key filename")
		nodeKeyHex    = flag.String("nodekeyhex", "", "private key as hex (for testing)")
		natdesc       = flag.String("nat", "none", "port mapping mechanism (any|none|upnp|pmp|extip:<IP>)")
		netrestrict   = flag.String("netrestrict", "", "restrict network communication to the given IP networks (CIDR masks)")
		runv5         = flag.Bool("v5", false, "run a v5 topic discovery bootnode")

		nodeKey *ecdsa.PrivateKey
		err     error
	)
	flag.Parse()
	getConfig(groupNum, groupNodesNum, listenAddr, nodeKeyFile)

	if *groupNodesNum == 0 {
		*groupNodesNum = 3
	}
	if *listenAddr == "" {
		*listenAddr = ":4440"
	}
	fmt.Printf("nodeKeyFile: %v, listenAddr: %v, group: %v, nodes: %v\n", *nodeKeyFile, *listenAddr, *groupNum, *groupNodesNum)
	natm, err := nat.Parse(*natdesc)
	if err != nil {
		fmt.Errorf("-nat: %v", err)
		return
	}
	switch {
	case *groupNodesNum == 0:
		fmt.Printf("Use -nodes must bigger than zero\n")
		return
	case *genKey != "":
		nodeKey, err = crypto.GenerateKey()
		if err != nil {
			//utils.Fatalf("could not generate key: %v", err)
		}
		if err = crypto.SaveECDSA(*genKey, nodeKey); err != nil {
			//utils.Fatalf("%v", err)
		}
		return
	case *nodeKeyFile == "" && *nodeKeyHex == "":
		fmt.Printf("Use -nodekey or -nodekeyhex to specify a private key\n")
		return
	case *nodeKeyFile != "" && *nodeKeyHex != "":
		fmt.Printf("Options -nodekey and -nodekeyhex are mutually exclusive\n")
		return
	case *nodeKeyFile != "":
		if nodeKey, err = crypto.LoadECDSA(*nodeKeyFile); err != nil {
			fmt.Printf("-nodekey: %v\n", err)
			return
		}
	case *nodeKeyHex != "":
		if nodeKey, err = crypto.HexToECDSA(*nodeKeyHex); err != nil {
			//utils.Fatalf("-nodekeyhex: %v", err)
		}
	}

	if *writeAddr {
		fmt.Printf("%v\n", discover.PubkeyID(&nodeKey.PublicKey))
		os.Exit(0)
	}

	var restrictList *netutil.Netlist
	if *netrestrict != "" {
		restrictList, err = netutil.ParseNetlist(*netrestrict)
		if err != nil {
			//utils.Fatalf("-netrestrict: %v", err)
		}
	}

	addr, err := net.ResolveUDPAddr("udp", *listenAddr)
	if err != nil {
		//utils.Fatalf("-ResolveUDPAddr: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		//utils.Fatalf("-ListenUDP: %v", err)
	}

	realaddr := conn.LocalAddr().(*net.UDPAddr)
	if natm != nil {
		if !realaddr.IP.IsLoopback() {
			go nat.Map(natm, nil, "udp", realaddr.Port, realaddr.Port, "ethereum discovery")
		}
		// TODO: react to external IP changes over time.
		if ext, err := natm.ExternalIP(); err == nil {
			realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
		}
	}

	if *runv5 {
		if _, err := discv5.ListenUDP(nodeKey, conn, realaddr, "", restrictList); err != nil {
			//utils.Fatalf("%v", err)
		}
	} else {
		cfg := discover.Config{
			PrivateKey:   nodeKey,
			AnnounceAddr: realaddr,
			NetRestrict:  restrictList,
		}
		if _, err := discover.ListenUDP(conn, cfg); err != nil {
			//utils.Fatalf("%v", err)
		}
	}

	// TODO: group
	fmt.Printf("groupNum: %v, groupNodesNum: %v\n", *groupNum, *groupNodesNum)
	if err := discover.InitGroup(int(*groupNum), int(*groupNodesNum)); err != nil {
	}

	select {}
}

type conf struct {
	Bootnode *bootnodeConf
}

type bootnodeConf struct {
	Nodekey string
	Addr    uint
	Group   uint
	Nodes   uint
}

func getConfig(groupNum, groupNodesNum *uint, listenAddr, nodeKeyFile *string) error {
	var cf conf
	var path string = "./conf.toml"
	if _, err := toml.DecodeFile(path, &cf); err != nil {
		//fmt.Printf("%v\n", err)
		return err
	}
	nkey := cf.Bootnode.Nodekey
	pt := cf.Bootnode.Addr
	gp := cf.Bootnode.Group
	ns := cf.Bootnode.Nodes
	if nkey != "" && *nodeKeyFile == "" {
		*nodeKeyFile = nkey
	}
	if pt != 0 && *listenAddr == "" {
		*listenAddr = fmt.Sprintf(":%v", pt)
	}
	if gp != 0 && *groupNum == 0 {
		*groupNum = gp
	}
	if ns != 0 && *groupNodesNum == 0 {
		*groupNodesNum = ns
	}
	return nil
}
