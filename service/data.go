package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"github.com/anyswap/fastmpc-service-middleware/db"
	"github.com/onrik/ethrpc"
)

func GetTestData() (string, error) {
	return db.Conn.GetStringValue("select name from test where id = ?", 1)
}

// GetGroupIdAndEnodes threshold 2/3, userAccountsAndIpPortAddr user1|ip:port user2 user3|ip:port
func GetGroupIdAndEnodes(threshold string, userAccountsAndIpPortAddr []string) (interface{}, error) {
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
		enodeData, _ := getJSONData(enodeRep)
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
	groupData, _ := getJSONData(groupRep)
	if err := json.Unmarshal(groupData, &groupJSON); err != nil {
		return nil, err
	}
	log.Info(fmt.Sprintf("\nGid = %s\n\n", groupJSON.Gid))

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	for _, v := range acct_enode {
		_, err = db.BatchExecute("insert into accounts_info(threshold, gid , user_account, ip_port, enode) values(?,?,?,?,?)", tx,
			threshold, groupJSON.Gid, v.Account, v.Ip_port, v.Enode)
		if err != nil {
			db.Conn.Rollback(tx)
			return nil, errors.New("internal db error " + err.Error())
		}
	}
	db.Conn.Commit(tx)

	return GroupIdAndEnodes{Gid: groupJSON.Gid, Enodes: enodeList}, nil
}

func getJSONData(successResponse json.RawMessage) ([]byte, error) {
	var rep response
	if err := json.Unmarshal(successResponse, &rep); err != nil {
		fmt.Println("getJSONData Unmarshal json fail:", err)
		return nil, err
	}
	if rep.Status != "Success" {
		return nil, errors.New(rep.Error)
	}
	repData, err := json.Marshal(rep.Data)
	if err != nil {
		fmt.Println("getJSONData Marshal json fail:", err)
		return nil, err
	}
	return repData, nil
}
