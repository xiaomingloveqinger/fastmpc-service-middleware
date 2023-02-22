package jobs

import (
	"encoding/hex"
	"encoding/json"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"github.com/anyswap/fastmpc-service-middleware/db"
	"github.com/onrik/ethrpc"
	"strings"
)

func listenKeyIdStatus() {
	list, err := db.Conn.GetStructValue("select ip_port, key_id, uuid from accounts_info where key_id is not null and status == 0 group by key_id, ip_port, uuid", Data{})
	if err != nil {
		log.Error("listenKeyIdStatus", "error", err.Error())
		return
	}
	for _, l := range list {
		var statusJSON reqAddrStatus
		d := l.(*Data)
		uuid := d.Uuid
		client := ethrpc.New("http://" + d.Ip_port)
		reqStatus, err := client.Call("smpc_getReqAddrStatus", d.Key_id)
		if err != nil {
			log.Error("smpc_getReqAddrStatus rpc error:" + err.Error())
			continue
		}
		statusJSONStr, err := common.GetJSONResult(reqStatus)
		if err != nil {
			log.Error("smpc_getReqAddrStatus=NotStart", "keyID", d.Key_id)
			continue
		}
		log.Info("smpc_getReqAddrStatus", "keyId", d.Key_id, "result", statusJSONStr)
		if err := json.Unmarshal([]byte(statusJSONStr), &statusJSON); err != nil {
			log.Error(err.Error())
			continue
		}
		if strings.ToLower(statusJSON.Status) != "pending" {
			log.Info("smpc_getReqAddrStatus", "smpc_getReqAddrStatus", statusJSON.Status, "keyID", d.Key_id)
			errMsg := statusJSON.Error
			tipMsg := statusJSON.Tip
			pub := statusJSON.PubKey
			stat, ok := ReqAddressStatusMap[strings.ToLower(statusJSON.Status)]
			if !ok {
				log.Error("can not find status in ReqAddressStatusMap")
				continue
			}
			tx, err := db.Conn.Begin()
			if err != nil {
				log.Error("internal db " + err.Error())
				continue
			}
			for _, reply := range statusJSON.AllReply {
				apprAcct := reply.Approver
				if apprAcct != "" {
					addr := ""
					if pub != "" {
						pubBuf, err := hex.DecodeString(pub)
						if err != nil {
							log.Error("invalid statusJson public key", "error", err.Error())
							continue
						}
						addr = common.PublicKeyBytesToAddress(pubBuf).String()
					}
					_, err := db.BatchExecute("update accounts_info set error = ? , tip = ? , reply_timestamp = ?, reply_status = ? , reply_initializer = ? , reply_enode = ? "+
						"mpc_address = ?, public_key = ? , status = ? where uuid = ? and user_account = ?", tx, errMsg, tipMsg, reply.TimeStamp, reply.Status, reply.Initiator, reply.Enode,
						addr, pub, stat, uuid, reply.Approver)
					if err != nil {
						db.Conn.Rollback(tx)
						log.Error("internal db error", "error", err.Error())
					}
				}
			}
			db.Conn.Commit(tx)
		}
	}
}

func init() {
	node.AddFunc("@every 30s", listenKeyIdStatus)
}