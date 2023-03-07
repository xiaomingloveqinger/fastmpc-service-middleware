package jobs

import (
	"encoding/hex"
	"encoding/json"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"github.com/anyswap/fastmpc-service-middleware/db"
	"github.com/onrik/ethrpc"
)

func listenSigningList() {
	signingAccts, err := db.Conn.GetStructValue("select user_account,ip_port from signs_detail where status = 0 group by user_account, ip_port", UserAccount{})
	if err != nil {
		log.Error("listenSigningList", "internal db error", err.Error())
		return
	}
	for _, acct := range signingAccts {
		a := acct.(*UserAccount)
		go runListenSigningList(a)
	}
}

func runListenSigningList(acct *UserAccount) {
	// get approve list of condominium account
	client := ethrpc.New("http://" + acct.Ip_port)
	log.Info("runListenSigningList " + acct.Ip_port)
	reqListRep, err := client.Call("smpc_getCurNodeSignInfo", acct.User_account)
	if err != nil {
		log.Error("runListenSigningList", "rpc call", err.Error())
		return
	}
	reqListJSON, _ := common.GetJSONData(reqListRep)
	log.Info("smpc_getCurNodeSignInfo", "msg", string(reqListJSON))
	var signing []signCurNodeInfo
	if err = json.Unmarshal(reqListJSON, &signing); err != nil {
		log.Error("Unmarshal signCurNodeInfo fail:", "msg", err.Error())
		return
	}
	singingKids, err := db.Conn.GetStructValue("select key_id from signing_list where user_account = ? and ip_port = ?", SigningKids{}, acct.User_account, acct.Ip_port)
	if err != nil {
		log.Error("internal db error " + err.Error())
		return
	}

	existed := extractSingingKids(singingKids)

	tx, err := db.Conn.Begin()
	if err != nil {
		log.Error("internal db error " + err.Error())
		return
	}
	for _, ing := range signing {
		delete(existed, ing.Key)
		c, err := db.Conn.GetIntValue("select count(key_id) from signing_list where key_id = ? and ip_port = ?", ing.Key, acct.Ip_port)
		if err != nil {
			log.Error("internal db error " + err.Error())
			db.Conn.Rollback(tx)
			return
		}
		if c > 0 {
			continue
		}
		pubBuf, err := hex.DecodeString(ing.PubKey)
		if err != nil {
			log.Error("invalid public key", "error", err.Error())
			return
		}
		addr := common.PublicKeyBytesToAddress(pubBuf).String()

		_, err = db.BatchExecute("insert into signing_list(user_account, group_id, key_id, key_type, `mode`, msg_context, msg_hash, nonce, public_key,mpc_address, threshold, `timestamp`, ip_port) values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
			tx, acct.User_account, ing.GroupID, ing.Key, ing.KeyType, ing.Mode, common.ConvertArrStrToStr(ing.MsgContext), common.ConvertArrStrToStr(ing.MsgHash), ing.Nonce, ing.PubKey, addr, ing.ThresHold, ing.TimeStamp, acct.Ip_port)
		if err != nil {
			log.Error("internal db error " + err.Error())
			db.Conn.Rollback(tx)
			return
		}
	}

	// update status if not exist
	if len(existed) > 0 {
		for k, _ := range existed {
			_, err = db.BatchExecute("update signing_list set status = 1 where key_id = ?", tx, k)
			log.Error("internal db error " + err.Error())
			db.Conn.Rollback(tx)
			return
		}
	}
	db.Conn.Commit(tx)
}

func extractSingingKids(d []interface{}) map[string]bool {
	m := make(map[string]bool)
	for _, v := range d {
		kids := v.(*SigningKids)
		m[kids.Key_id] = true
	}
	return m
}

func init() {
	jobs.AddFunc("@every 1m", listenSigningList)
}
