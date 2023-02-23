package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"github.com/anyswap/fastmpc-service-middleware/db"
	common2 "github.com/anyswap/fastmpc-service-middleware/internal/common"
	"github.com/google/uuid"
	"github.com/onrik/ethrpc"
	"strings"
)

func doKeyGen(rsv string, msg string) (interface{}, error) {
	err := common.VerifyAccount(rsv, msg)
	if err != nil {
		return nil, err
	}
	req := TxDataReqAddr{}
	err = json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return nil, err
	}
	if req.Mode != "3" {
		return nil, errors.New("service keygen mod must be 3")
	}
	req.Account = strings.ToLower(req.Account)
	ipStr, err := db.Conn.GetStringValue("select ip_port from accounts_info where gid = ? and user_account = ? and threshold = ? and key_id is null and uuid = ?",
		req.GroupID, req.Account, req.ThresHold, req.Uuid)
	if err != nil {
		return nil, errors.New("GetStringValue error " + err.Error())
	}
	if ipStr == "" {
		return nil, errors.New("Can not find ip port through uuid " + req.Uuid)
	}
	client := ethrpc.New("http://" + ipStr)
	reqKeyID, err := client.Call("smpc_reqKeyGen", rsv, msg)
	if err != nil {
		return nil, errors.New("smpc_reqKeyGen error " + err.Error())
	}
	keyID, err := common.GetJSONResult(reqKeyID)
	if err != nil {
		return nil, errors.New("getJSONResult error" + err.Error())
	}
	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	_, err = db.BatchExecute("update accounts_info set key_id = ? where uuid = ?",
		tx, keyID, req.Uuid)
	if err != nil {
		db.Conn.Rollback(tx)
		return nil, errors.New("internal db error " + err.Error())
	}
	_, err = db.BatchExecute("insert into groups_info(tx_type, account, nonce, key_type, group_id, thres_hold, mode, accept_timeout, sigs , key_id, uuid, timestamp) "+
		"values(?,?,?,?,?,?,?,?,?,?,?,?)", tx, req.TxType, req.Account, req.Nonce, req.Keytype, req.GroupID, req.ThresHold, req.Mode, req.AcceptTimeOut, req.Sigs, keyID, req.Uuid, req.TimeStamp)
	if err != nil {
		db.Conn.Rollback(tx)
		return nil, errors.New("internal db error " + err.Error())
	}
	if err = db.Conn.Commit(tx); err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	return keyID, nil
}

func getGroupIdAndEnodesByRawData(raw string) (interface{}, error) {
	type Msg struct {
		Threshold                 string
		UserAccountsAndIpPortAddr []string
	}
	m := Msg{}
	err := json.Unmarshal(common2.FromHex(raw), &m)
	if err != nil {
		return nil, err
	}
	return getGroupId(m.Threshold, m.UserAccountsAndIpPortAddr)
}

// getGroupId threshold 2/3, userAccountsAndIpPortAddr user1|ip:port user2 user3|ip:port
func getGroupId(threshold string, userAccountsAndIpPortAddr []string) (interface{}, error) {
	_, p2, err := common.CheckThreshold(threshold)
	if err != nil {
		return nil, err
	}
	accounts, ipPorts, err := common.CheckUserAccountsAndIpPortAddr(userAccountsAndIpPortAddr)
	if err != nil {
		return nil, err
	}
	if p2 != len(accounts) {
		return nil, errors.New("threshold value and accounts number can not match")
	}
	var selectCount int
	var filledIpPort string
	for _, v := range ipPorts {
		if v == "" {
			selectCount++
		} else {
			filledIpPort = filledIpPort + v + ","
		}
	}
	filledIpPort = filledIpPort[0 : len(filledIpPort)-1]
	filledIpPortCount, err := db.Conn.GetIntValue("select count(*) from nodes_info where ip_addr in (?)", filledIpPort)
	if err != nil {
		return nil, err
	}
	if filledIpPortCount != len(ipPorts)-selectCount {
		return nil, errors.New("filled ip port not valid")
	}
	type ip struct {
		Ip_addr string
	}
	totalNodes, err := db.Conn.GetIntValue("select count(ip_addr) from nodes_info")
	if err != nil {
		return nil, err
	}
	if totalNodes < selectCount {
		return nil, errors.New("total nodes less than needed")
	}
	allIp, err := db.Conn.GetStructValue("select ip_addr from nodes_info where ip_addr not in (?) ORDER BY RAND() limit ?", ip{}, filledIpPort, selectCount)
	if err != nil {
		return nil, err
	}
	index := 0
	var acct_enode []account_enode
	var enodeList []string
	for i, acc := range accounts {
		var enode string
		var ipStr string
		if ipPorts[i] == "" {
			ipStr = allIp[index].(*ip).Ip_addr
			index++
		} else {
			ipStr = ipPorts[i]
		}
		client := ethrpc.New("http://" + ipStr)
		enodeRep, err := client.Call("smpc_getEnode")
		if err != nil {
			return nil, errors.New("IP_addr" + allIp[index].(ip).Ip_addr + " can not reach")
		}
		log.Info(fmt.Sprintf("getEnode = %s\n\n", enodeRep))
		type dataEnode struct {
			Enode string `json:"Enode"`
		}
		var enodeJSON dataEnode
		enodeData, _ := common.GetJSONData(enodeRep)
		if err := json.Unmarshal(enodeData, &enodeJSON); err != nil {
			return nil, err
		}
		log.Info(fmt.Sprintf("enode = %s\n", enodeJSON.Enode))
		enode = enodeJSON.Enode

		acct_enode = append(acct_enode, account_enode{
			Enode:   enode,
			Account: acc,
			Ip_port: ipStr,
		})

		enodeList = append(enodeList, enode)
	}

	client := ethrpc.New("http://" + acct_enode[0].Ip_port)
	// get gid by send createGroup
	groupRep, err := client.Call("smpc_createGroup", threshold, enodeList)
	if err != nil {
		return nil, err
	}
	log.Info(fmt.Sprintf("smpc_createGroup = %s\n", groupRep))
	var groupJSON groupInfo
	groupData, _ := common.GetJSONData(groupRep)
	if err := json.Unmarshal(groupData, &groupJSON); err != nil {
		return nil, err
	}
	log.Info(fmt.Sprintf("\nGid = %s\n\n", groupJSON.Gid))

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	uid := uuid.New().String()
	sigs := ""
	for _, v := range acct_enode {
		_, err = db.BatchExecute("insert into accounts_info(threshold, gid , user_account, ip_port, enode, uuid) values(?,?,?,?,?,?)", tx,
			threshold, groupJSON.Gid, strings.ToLower(v.Account), v.Ip_port, v.Enode, uid)
		if err != nil {
			db.Conn.Rollback(tx)
			return nil, errors.New("internal db error " + err.Error())
		}
		sigs += common.StripEnode(v.Enode) + ":" + strings.ToLower(v.Account) + "|"
	}
	if err = db.Conn.Commit(tx); err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	return Group{Gid: groupJSON.Gid, Sigs: sigs, Uuid: uid}, nil
}
