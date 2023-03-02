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

//listenSignKidStatus listen sign keyid status and stored it into db
func listenSignKidStatus() {
	list, err := db.Conn.GetStructValue("select key_id,pubkey from signs_info where key_id is not null and status = 0", SignKids{})
	if err != nil {
		log.Error("listenKeygenKidStatus", "error", err.Error())
		return
	}
	for _, l := range list {
		d := l.(*SignKids)
		kid := d.Key_id
		ps, err := db.Conn.GetStructValue("select ip_port from accounts_info where public_key = ? and status = 1", IpData{}, d.Pubkey)
		if err != nil {
			log.Error("internal db " + err.Error())
			return
		}
		for _, p := range ps {
			ip := p.(*IpData)
			client := ethrpc.New("http://" + ip.Ip_port)
			reqStatus, err := client.Call("smpc_getSignStatus", kid)
			if err != nil {
				log.Error("smpc_getReqAddrStatus rpc error:" + err.Error())
				continue
			}
			statusJSONStr, err := common.GetJSONResult(reqStatus)
			if err != nil {
				log.Error("smpc_getReqAddrStatus=NotStart", "keyID", d.Key_id, "error", err.Error())
				continue
			}
			var statusJSON SignStatus
			log.Info("smpc_getReqAddrStatus", "keyId", d.Key_id, "result", statusJSONStr)
			if err := json.Unmarshal([]byte(statusJSONStr), &statusJSON); err != nil {
				log.Error(err.Error())
				return
			}
			if strings.ToLower(statusJSON.Status) != "pending" {
				log.Info("smpc_getSignStatus", "smpc_getSignStatus", statusJSON.Status, "keyID", d.Key_id)
				pub := statusJSON.PubKey
				stat, ok := ReqAddressStatusMap[strings.ToLower(statusJSON.Status)]
				if !ok {
					log.Error("can not find status in ReqAddressStatusMap")
					return
				}
				tx, err := db.Conn.Begin()
				if err != nil {
					log.Error("internal db " + err.Error())
					return
				}
				for _, reply := range statusJSON.AllReply {
					addr := ""
					if pub != "" {
						pubBuf, err := hex.DecodeString(pub)
						if err != nil {
							log.Error("invalid statusJson public key", "error", err.Error())
							return
						}
						addr = common.PublicKeyBytesToAddress(pubBuf).String()
					}
					if reply.Approver == "" {
						log.Warn("reply approver can not be blank " + kid)
						return
					}
					exist, err := db.Conn.GetIntValue("select count(kid) from signs_detail where key_id = ? and user_account = ?", kid, reply.Approver)
					if err != nil {
						log.Error("internal db error", "error", err.Error())
						return
					}
					if exist > 0 {
						log.Info("already exist record ", "keyid", kid, "user_account", reply.Approver)
						continue
					}
					_, err = db.BatchExecute("insert into signs_detail(user_account, group_id, threshold, msg_hash, msg_context, rsv, public_key, mpc_address, reply_initializer, reply_status, reply_timestamp, reply_enode, initiator_public_key, key_type, mode, error, tip, status) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", tx,
						reply.Approver, statusJSON.GroupID, statusJSON.ThresHold, common.ConvertArrStrToStr(statusJSON.MsgHash), common.ConvertArrStrToStr(statusJSON.MsgContext),
						common.ConvertArrStrToStr(statusJSON.Rsv), pub, addr, reply.Initiator, reply.Status, reply.TimeStamp, reply.Enode, statusJSON.Initiator, statusJSON.Keytype, statusJSON.Mode, statusJSON.Error, statusJSON.Tip, stat)
					if err != nil {
						db.Conn.Rollback(tx)
						log.Error("internal db error", "error", err.Error())
						return
					}
				}
				_, err = db.BatchExecute("update signs_info set status = 1 where kid = ?", tx, kid)
				if err != nil {
					db.Conn.Rollback(tx)
					log.Error("internal db error", "error", err.Error())
					return
				}
				db.Conn.Commit(tx)
			}
		}
	}
}

func init() {
	jobs.AddFunc("@every 1m", listenSignKidStatus)
}
