/*
 *  Copyright (C) 2018-2019  Fusion Foundation Ltd. All rights reserved.
 *  Copyright (C) 2018-2019  huangweijun@fusion.org
 *
 *  This library is free software; you can redistribute it and/or
 *  modify it under the Apache License, Version 2.0.
 *
 *  This library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package layer2

import (
	"context"
	"errors"

	//"math"
	"fmt"
	"net"
	"sort"
	"time"

	"github.com/EricBui0512/dcrm-walletService/crypto"
	"github.com/EricBui0512/dcrm-walletService/p2p"
	"github.com/EricBui0512/dcrm-walletService/p2p/discover"
	"github.com/EricBui0512/dcrm-walletService/rpc"
)

// txs start
func DcrmProtocol_sendToGroupOneNode(msg string) (string, error) {
	return discover.SendToGroup(discover.NodeID{}, msg, false, DcrmProtocol_type, nil)
}

// broadcast
// to group's nodes
func DcrmProtocol_broadcastInGroupAll(msg string) { // within self
	BroadcastToGroup(discover.NodeID{}, msg, DcrmProtocol_type, true)
}

func DcrmProtocol_broadcastInGroupOthers(msg string) { // without self
	BroadcastToGroup(discover.NodeID{}, msg, DcrmProtocol_type, false)
}

// unicast
// to anyone
func DcrmProtocol_sendMsgToNode(toid discover.NodeID, toaddr *net.UDPAddr, msg string) error {
	fmt.Printf("==== SendMsgToNode() ====\n")
	return discover.SendMsgToNode(toid, toaddr, msg)
}

// to peers
func DcrmProtocol_sendMsgToPeer(enode string, msg string) error {
	return SendMsgToPeer(enode, msg)
}

// callback
// receive private key
func DcrmProtocol_registerPriKeyCallback(recvPrivkeyFunc func(interface{})) {
	discover.RegisterPriKeyCallback(recvPrivkeyFunc)
}

func Sdk_callEvent(msg string, fromID string) {
	fmt.Printf("Sdk_callEvent\n")
	Sdk_callback(msg, fromID)
}

// receive message form peers
func DcrmProtocol_registerRecvCallback(recvDcrmFunc func(interface{}) <-chan string) {
	Dcrm_callback = recvDcrmFunc
}
func Dcrm_callEvent(msg string) {
	Dcrm_callback(msg)
}

// receive message from dccp
func DcrmProtocol_registerMsgRecvCallback(dcrmcallback func(interface{}) <-chan string) {
	discover.RegisterDcrmMsgCallback(dcrmcallback)
}

// receive message from dccp result
func DcrmProtocol_registerMsgRetCallback(dcrmcallback func(interface{})) {
	discover.RegisterDcrmMsgRetCallback(dcrmcallback)
}

// get info
func DcrmProtocol_getGroup() (int, string) {
	return getGroup(discover.NodeID{}, DcrmProtocol_type)
}

func (dcrm *DcrmAPI) Version(ctx context.Context) (v string) {
	return ProtocolVersionStr
}
func (dcrm *DcrmAPI) Peers(ctx context.Context) []*p2p.PeerInfo {
	var ps []*p2p.PeerInfo
	for _, p := range dcrm.dcrm.peers {
		ps = append(ps, p.peer.Info())
	}

	return ps
}

// Protocols returns the whisper sub-protocols ran by this particular client.
func (dcrm *Dcrm) Protocols() []p2p.Protocol {
	return []p2p.Protocol{dcrm.protocol}
}

// p2p layer 2
// New creates a Whisper client ready to communicate through the Ethereum P2P network.
func DcrmNew(cfg *Config) *Dcrm {
	dcrm := &Dcrm{
		peers: make(map[discover.NodeID]*peer),
		quit:  make(chan struct{}),
		cfg:   cfg,
	}

	// p2p dcrm sub protocol handler
	dcrm.protocol = p2p.Protocol{
		Name:    ProtocolName,
		Version: ProtocolVersion,
		Length:  NumberOfMessageCodes,
		Run:     HandlePeer,
		NodeInfo: func() interface{} {
			return map[string]interface{}{
				"version": ProtocolVersionStr,
			}
		},
		PeerInfo: func(id discover.NodeID) interface{} {
			if p := emitter.peers[id]; p != nil {
				return p.peerInfo
			}
			return nil
		},
	}

	return dcrm
}

// TODO callback
func recvPrivkeyInfo(msg interface{}) {
	fmt.Printf("==== recvPrivkeyInfo() ====\n")
	fmt.Println("recvprikey,msg = ", msg)
	//TODO
	//store privatekey slice
	time.Sleep(time.Duration(10) * time.Second)
	DcrmProtocol_broadcastInGroupOthers("aaaa")
}

// other
// Start implements node.Service, starting the background data propagation thread
// of the Whisper protocol.
func (dcrm *Dcrm) Start(server *p2p.Server) error {
	return nil
}

// Stop implements node.Service, stopping the background data propagation thread
// of the Whisper protocol.
func (dcrm *Dcrm) Stop() error {
	return nil
}

// APIs returns the RPC descriptors the Whisper implementation offers
func (dcrm *Dcrm) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: ProtocolName,
			Version:   ProtocolVersionStr,
			Service:   &DcrmAPI{dcrm: dcrm},
			Public:    true,
		},
	}
}

func DcrmProtocol_getEnodes() (int, string) {
	return getGroup(discover.NodeID{}, DcrmProtocol_type)
}

func RegisterUpdateOrderCacheCallback(recvDcrmFunc func(interface{})) {
	discover.RegisterUpdateOrderCacheCallback(recvDcrmFunc)
}

// =============================== DCRM =================================
func SendMsg(msg string) {
	//BroadcastToGroup(discover.NodeID{}, msg, DcrmProtocol_type)
	DcrmProtocol_broadcastInGroupOthers(msg)
}

func SendToDcrmGroupAllNodes(msg string) (string, error) {
	return discover.SendToGroup(discover.NodeID{}, msg, true, DcrmProtocol_type, nil)
}

func RegisterRecvCallback(recvPrivkeyFunc func(interface{})) {
	discover.RegisterPriKeyCallback(recvPrivkeyFunc)
}

func RegisterDcrmCallback(dcrmcallback func(interface{}) <-chan string) {
	discover.RegisterDcrmMsgCallback(dcrmcallback)
	DcrmProtocol_registerRecvCallback(dcrmcallback)
}

func RegisterDcrmRetCallback(dcrmcallback func(interface{})) {
	discover.RegisterDcrmMsgRetCallback(dcrmcallback)
}

func GetGroup() (int, string) {
	return DcrmProtocol_getGroup()
}

func GetEnodes() (int, string) {
	return GetGroup()
}

func RegisterSendCallback(callbackfunc func(interface{})) {
	discover.RegisterSendCallback(callbackfunc)
}

func ParseNodeID(enode string) string {
	node, _ := discover.ParseNode(enode)
	return node.ID.String()
}

// ================   API   SDK    =====================
func SdkProtocol_sendToGroupOneNode(gID, msg string) (string, error) {
	gid, _ := discover.HexID(gID)
	if checkExistGroup(gid) == false {
		fmt.Printf("sendToGroupOneNode, group gid: %v not exist\n", gid)
		return "", errors.New("sendToGroupOneNode, gid not exist")
	}
	g := getSDKGroupNodes(gid)
	return discover.SendToGroup(gid, msg, false, Sdkprotocol_type, g)
}

func getSDKGroupNodes(gid discover.NodeID) []*discover.Node {
	g := make([]*discover.Node, 0)
	_, xvcGroup := getGroupSDK(gid)
	if xvcGroup == nil {
		return g
	}
	for _, rn := range xvcGroup.Nodes {
		n := discover.NewNode(rn.ID, rn.IP, rn.UDP, rn.TCP)
		g = append(g, n)
	}
	return g
}

func SdkProtocol_SendToGroupAllNodes(gID, msg string) (string, error) {
	gid, _ := discover.HexID(gID)
	if checkExistGroup(gid) == false {
		e := fmt.Sprintf("SendGroupAllNodes, group gid: %v not exist", gid)
		fmt.Println(e)
		return "", errors.New(e)
	}
	g := getSDKGroupNodes(gid)
	return discover.SendToGroup(gid, msg, true, Sdkprotocol_type, g)
}

func SdkProtocol_broadcastInGroupOthers(gID, msg string) (string, error) { // without self
	gid, _ := discover.HexID(gID)
	if checkExistGroup(gid) == false {
		e := fmt.Sprintf("broadcastInGroupOthers, group gid: %v not exist", gid)
		fmt.Println(e)
		return "", errors.New(e)
	}
	return BroadcastToGroup(gid, msg, Sdkprotocol_type, false)
}

func SdkProtocol_broadcastInGroupAll(gID, msg string) (string, error) { // within self
	gid, _ := discover.HexID(gID)
	if checkExistGroup(gid) == false {
		e := fmt.Sprintf("broadcastInGroupAll, group gid: %v not exist", gid)
		fmt.Println(e)
		return "", errors.New(e)
	}
	return BroadcastToGroup(gid, msg, Sdkprotocol_type, true)
}

func SdkProtocol_getGroup(gID string) (int, string) {
	//fmt.Printf("getGroup, gid: %v\n", gID)
	gid, _ := discover.HexID(gID)
	if checkExistGroup(gid) == false {
		fmt.Printf("broadcastInGroupAll, group gid: %v not exist\n", gid)
		return 0, ""
	}
	return getGroup(gid, Sdkprotocol_type)
}

func checkExistGroup(gid discover.NodeID) bool {
	if SdkGroup[gid] != nil {
		if SdkGroup[gid].Type == "1+2" || SdkGroup[gid].Type == "1+1+1" {
			return true
		}
	}
	return false
}

//	---------------------   API  callback   ----------------------
//
// recv from broadcastInGroup...
func SdkProtocol_registerBroadcastInGroupCallback(recvSdkFunc func(interface{}, string)) {
	Sdk_callback = recvSdkFunc
}

// recv from sendToGroup...
func SdkProtocol_registerSendToGroupCallback(sdkcallback func(interface{}, string) <-chan string) {
	discover.RegisterSdkMsgCallback(sdkcallback)
}

// recv return from sendToGroup...
func SdkProtocol_registerSendToGroupReturnCallback(sdkcallback func(interface{}, string)) {
	discover.RegisterSdkMsgRetCallback(sdkcallback)
}

// 1 + 1 + 1
func CreateSDKGroup(mode string, enodes []string) (string, int, string) {
	count := len(enodes)
	sort.Sort(sort.StringSlice(enodes))
	enode := []*discover.Node{}
	id := []byte("")
	for _, un := range enodes {
		fmt.Printf("for un: %v\n", un)
		node, err := discover.ParseNode(un)
		if err != nil {
			fmt.Printf("CreateSDKGroup, parse err: %v\n", un)
			return "", 0, "enode wrong format"
		}
		fmt.Printf("for selfid: %v, node.ID: %v\n", selfid, node.ID)
		if selfid != node.ID {
			p := emitter.peers[node.ID]
			if p == nil {
				fmt.Printf("CreateSDKGroup, peers err: %v\n", un)
				return "", 0, "enode is not peer"
			}
		}
		n := fmt.Sprintf("%v", node.ID)
		fmt.Printf("CreateSDKGroup, n: %v\n", n)
		if len(id) == 0 {
			id = crypto.Keccak512([]byte(node.ID.String()))
		} else {
			id = crypto.Keccak512(id, []byte(node.ID.String()))
		}
		enode = append(enode, node)
	}
	gid, err := discover.BytesID(id)
	fmt.Printf("CreateSDKGroup, gid <- id: %v, err: %v\n", gid, err)
	for i, g := range SdkGroup {
		if i == gid {
			return gid.String(), len(g.Nodes), "group is exist"
		}
	}
	retErr := discover.StartCreateSDKGroup(gid, mode, enode, "1+1+1")
	return gid.String(), count, retErr
}

func GetEnodeStatus(enode string) (string, string) {
	return discover.GetEnodeStatus(enode)
}

func CheckAddPeer(enodes []string) error {
	addpeer := false
	var nodes []*discover.Node
	for _, enode := range enodes {
		node, err := discover.ParseNode(enode)
		if err != nil {
			msg := fmt.Sprintf("CheckAddPeer, parse err enode: %v", enode)
			return errors.New(msg)
		}
		if selfid == node.ID {
			continue
		}

		status, _ := GetEnodeStatus(enode)
		if status == "OffLine" {
			msg := fmt.Sprintf("CheckAddPeer, enode: %v offline", enode)
			return errors.New(msg)
		}
		p := emitter.peers[node.ID]
		if p == nil {
			addpeer = true
			nodes = append(nodes, node)
			go p2pServer.AddPeer(node)
		}
	}
	if addpeer {
		fmt.Printf("CheckAddPeer, waitting add peer ...\n")
		count := 0
		for _, node := range nodes {
			p := emitter.peers[node.ID]
			if p == nil {
				time.Sleep(time.Duration(1) * time.Second)
				count += 1
				if count > (len(nodes) * 10) {
					fmt.Printf("CheckAddPeer, add peer failed node: %v\n", node)
					msg := fmt.Sprintf("CheckAddPeer, add peer failed node: %v", node)
					return errors.New(msg)
				}
				continue
			}
			fmt.Printf("CheckAddPeer, add peer success node: %v\n", node)
		}
	}
	return nil
}
