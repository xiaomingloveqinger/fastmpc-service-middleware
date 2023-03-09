package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/chains/types"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"github.com/anyswap/fastmpc-service-middleware/db"
	common2 "github.com/anyswap/fastmpc-service-middleware/internal/common"
	"github.com/google/uuid"
	"github.com/onrik/ethrpc"
	"strconv"
	"strings"
)

func acceptSign(rsv string, msg string) (interface{}, error) {
	err := common.VerifyAccount(rsv, msg)
	if err != nil {
		return nil, err
	}
	req := AcceptSignData{}
	err = json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return nil, err
	}
	if len(req.MsgHash) != len(req.MsgContext) {
		return nil, errors.New("message hash and message context length not match")
	}

	if len(req.MsgHash) == 0 {
		return nil, errors.New("message hash and message context can not be blank")
	}

	if common.IsSomeOneBlank(req.Accept, req.Nonce, req.Key, req.TxType, req.TimeStamp) {
		return nil, errors.New("request param not valid")
	}

	if req.TxType != "ACCEPTSIGN" {
		return nil, errors.New("invalid tx type")
	}

	if req.Accept != "AGREE" && req.Accept != "DISAGREE" {
		return nil, errors.New("invalid accept value")
	}

	switch types.ChainType(req.ChainType) {
	case types.EVM:
		for i, hash := range req.MsgHash {
			if !types.EvmChain.ValidateUnsignedTransactionHash(req.MsgContext[i], hash) {
				return nil, errors.New("message hash and msg context value not match")
			}
		}
	default:
		return nil, errors.New("unrecognized chain type")
	}

	ipPort, err := db.Conn.GetStringValue("select ip_port from signs_detail where key_id = ? and user_account = ?", req.Key, strings.ToLower(req.Account))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	if common.IsBlank(ipPort) {
		return nil, errors.New("invalid request param can not find ip port")
	}

	client := ethrpc.New("http://" + ipPort)
	// send rawTx
	acceptSignRep, err := client.Call("smpc_acceptSigning", rsv, msg)
	if err != nil {
		return nil, err
	}
	// get result
	acceptRet, err := common.GetJSONResult(acceptSignRep)
	if err != nil {
		return nil, err
	}
	log.Info("smpc_acceptSign: ", "result", acceptRet)

	_, err = db.Conn.CommitOneRow("insert into signs_result(key_type,account,nonce,key_id,msg_hash,msg_context,timestamp,accept) values(?,?,?,?,?,?,?,?)",
		req.TxType, req.Account, req.Nonce, req.Key, common.ConvertArrStrToStr(req.MsgHash), common.ConvertArrStrToStr(req.MsgContext), req.TimeStamp, req.Accept)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	return acceptRet, nil
}

func getSignHistory(userAccount string) (interface{}, error) {
	if !common.CheckEthereumAddress(userAccount) {
		return nil, errors.New("user account is not valid")
	}
	l, err := db.Conn.GetStructValue("select * from signs_detail where user_account = ?", SignHistory{}, userAccount)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	return l, nil
}

func getApprovalList(userAccount string) (interface{}, error) {
	if !common.CheckEthereumAddress(userAccount) {
		return nil, errors.New("user account is not valid")
	}
	l, err := db.Conn.GetStructValue("select * from signing_list where user_account = ?", SignCurNodeInfo{}, userAccount)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	return l, nil
}

func getUnsigedTransactionHash(unsignedTx string, chain int) (interface{}, error) {
	var c types.Chain
	switch types.ChainType(chain) {
	case types.EVM:
		c = types.EvmChain
	default:
		return nil, errors.New("unrecognized chain")
	}
	hash, err := c.GetUnsignedTransactionHash(unsignedTx)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func doSign(rsv string, msg string) (interface{}, error) {
	err := common.VerifyAccount(rsv, msg)
	if err != nil {
		return nil, err
	}
	req := TxDataSign{}
	err = json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return nil, err
	}
	if len(req.MsgHash) != len(req.MsgContext) {
		return nil, errors.New("message hash and message context length not match")
	}

	if len(req.MsgHash) == 0 {
		return nil, errors.New("message hash and message context can not be blank")
	}

	switch types.ChainType(req.ChainType) {
	case types.EVM:
		for i, hash := range req.MsgHash {
			if !types.EvmChain.ValidateUnsignedTransactionHash(req.MsgContext[i], hash) {
				return nil, errors.New("message hash and msg context value not match")
			}
		}
	default:
		return nil, errors.New("unrecognized chain type")
	}

	ipPort, err := db.Conn.GetStringValue("select ip_port from accounts_info where public_key = ? and user_account = ? and status = 1", req.PubKey, strings.ToLower(req.Account))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	if ipPort == "" {
		return nil, errors.New("can not find pub key and account responded ip port")
	}
	client := ethrpc.New("http://" + ipPort)
	reqKeyID, err := client.Call("smpc_signing", rsv, msg)
	if err != nil {
		return nil, err
	}
	keyID, err := common.GetJSONResult(reqKeyID)
	if err != nil {
		return nil, err
	}
	log.Info("smpc_sign keyID = %s", keyID)

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, errors.New("internal db error" + err.Error())
	}
	_, err = db.BatchExecute("insert into signs_info(account,nonce,pubkey,msg_hash,msg_context,key_type,group_id,threshold,`mod`,accept_timeout,`timestamp`,key_id) values(?,?,?,?,?,?,?,?,?,?,?,?)",
		tx, req.Account, req.Nonce, req.PubKey, common.ConvertArrStrToStr(req.MsgHash), common.ConvertArrStrToStr(req.MsgContext), req.Keytype, req.GroupID, req.ThresHold, req.Mode, req.AcceptTimeOut, req.TimeStamp, keyID)
	if err != nil {
		db.Conn.Rollback(tx)
		return nil, errors.New("internal db error" + err.Error())
	}
	pubBuf, err := hex.DecodeString(req.PubKey)
	if err != nil {
		return nil, errors.New("invalid req pub key")
	}
	addr := common.PublicKeyBytesToAddress(pubBuf).String()
	accts, err := db.Conn.GetStructValue("select user_account, enode, ip_port from accounts_info where public_key = ?", Account{}, req.PubKey)
	if err != nil {
		db.Conn.Rollback(tx)
		return nil, errors.New("internal db error" + err.Error())
	}
	if len(accts) == 0 {
		db.Conn.Rollback(tx)
		return nil, errors.New("invalid public key")
	}
	for _, acct := range accts {
		a := acct.(*Account)
		_, err = db.BatchExecute("insert into signs_detail(key_id, user_account, group_id, threshold, msg_hash, msg_context, public_key, mpc_address, key_type, mode, status, enode, ip_port) values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
			tx, keyID, a.User_account, req.GroupID, req.ThresHold, common.ConvertArrStrToStr(req.MsgHash), common.ConvertArrStrToStr(req.MsgContext), req.PubKey, addr, req.Keytype, req.Mode, 0, a.Enode, a.Ip_port)
		if err != nil {
			db.Conn.Rollback(tx)
			return nil, errors.New("internal db error" + err.Error())
		}
	}

	db.Conn.Commit(tx)
	return keyID, nil
}

func getAccountList(userAccount string) (interface{}, error) {
	if !common.CheckEthereumAddress(userAccount) {
		return nil, errors.New("invalid userAccount")
	}

	l, err := db.Conn.GetStructValue("select * from accounts_info where user_account = ? and status = 1 ", RespAddr{}, userAccount)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func getReqAddrStatus(keyId string) (interface{}, error) {
	if !common.ValidateKeyId(keyId) {
		return nil, errors.New("keyId is not valid")
	}

	v, err := db.Conn.GetStructValue("select status, user_account, key_id, public_key, mpc_address, initializer, reply_status ,reply_timestamp ,reply_enode, gid , threshold from accounts_info where key_id = ?", RespAddr{}, keyId)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	if len(v) == 0 {
		return nil, errors.New("no such keyId")
	}

	return v, nil
}

func doKeyGenByRawData(raw string) (interface{}, error) {
	type Msg struct {
		Rsv string
		Msg string
	}
	m := Msg{}
	err := json.Unmarshal(common2.FromHex(raw), &m)
	if err != nil {
		return nil, err
	}
	return doKeyGen(m.Rsv, m.Msg)
}

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
	if req.Mode != "2" {
		return nil, errors.New("service keygen mod must be 2")
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

func getGroupIdByRawData(raw string) (interface{}, error) {
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
	if filledIpPort != "" {
		filledIpPort = filledIpPort[0 : len(filledIpPort)-1]
	}
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
		sigs += common.StripEnode(v.Enode) + ":" + strings.ToLower(v.Account) + ":"
	}
	if err = db.Conn.Commit(tx); err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	return Group{Gid: groupJSON.Gid, Sigs: strconv.Itoa(len(acct_enode)) + ":" + sigs[:len(sigs)-1], Uuid: uid}, nil
}
