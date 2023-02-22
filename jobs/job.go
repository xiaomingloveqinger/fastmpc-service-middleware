package jobs

import "github.com/robfig/cron/v3"

var node MpcNodesInfo

type MpcNodesInfo struct {
	*cron.Cron
}

func init() {
	node = MpcNodesInfo{
		Cron: cron.New(),
	}
	node.Start()
}
