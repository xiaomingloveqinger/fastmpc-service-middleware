package jobs

import (
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/db"
	"github.com/robfig/cron/v3"
)

var node MpcNodesInfo

type MpcNodesInfo struct {
	*cron.Cron
}

type Node struct {
	Ip_addr     string
	Name        string
	Email       string
	Telegram_id string
	Enode       string
}

func getRegisteredNodeInfo() {
	log.Info("listen register nodes")
	//TODO: tmp static data will be removed in the future
	var nodes []Node
	nodes = append(nodes, Node{
		Ip_addr:     "127.0.0.1:3794",
		Name:        "test node1",
		Email:       "test1@gmail.com",
		Telegram_id: "test1",
		Enode:       "enode://748ba7475b0da18887480871eb6a41c0b207c2056bf9e0cbe2d25677fef9849e3ec82d038e3d820ba9586abd1a1327555c63c34b71d9b8bccd7a1e3bedeca47b@127.0.0.1:30823",
	})

	nodes = append(nodes, Node{
		Ip_addr:     "127.0.0.1:3793",
		Name:        "test node2",
		Email:       "test2@gmail.com",
		Telegram_id: "test2",
		Enode:       "enode://2e2b74160a62114e8901668022ab8df0d30ae9c69a48100ab70d50da4713ca6d71ca1bee30bd60a505077ea1c1c2b67b423ed75d535599c3be2b46f397de1a96@127.0.0.1:30824",
	})

	nodes = append(nodes, Node{
		Ip_addr:     "127.0.0.1:3792",
		Name:        "test node3",
		Email:       "test3@gmail.com",
		Telegram_id: "test3",
		Enode:       "enode://08ba43b0715bb27e03911592d3fed49a22f49ecaf1628d44c4c7d4a8914b86d423716d7316d98a80f4d280b852a33ba7972f0a744e386375e6d58a12fab96752@127.0.0.1:30825",
	})

	for _, n := range nodes {
		c, err := db.Conn.GetIntValue("select count(ip_addr) from nodes_info where ip_addr = ?", n.Ip_addr)
		if err != nil {
			log.Error("DB Error", "GetIntValue", err.Error())
			return
		}
		if c > 0 {
			log.Info("DB Exist", "Ip addr", n.Ip_addr)
			continue
		}
		_, err = db.Conn.CommitOneRow("insert into nodes_info(ip_addr,name,email,telegram_id,enode) values(?,?,?,?,?)", n.Ip_addr, n.Name, n.Email, n.Telegram_id, n.Enode)
		if err != nil {
			log.Error("DB Error", "CommitOneRow", err.Error())
			return
		}
	}
}

func init() {
	node = MpcNodesInfo{
		Cron: cron.New(),
	}
	node.AddFunc("@every 30s", getRegisteredNodeInfo)
	node.Start()
}
