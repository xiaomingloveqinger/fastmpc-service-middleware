package service

type account_enode struct {
	Account string `json:"account"`
	Enode   string `json:"enode"`
	Ip_port string `json:"ip_port"`
}

type groupInfo struct {
	Gid    string      `json:"Gid"`
	Mode   string      `json:"Mode"`
	Count  int         `json:"Count"`
	Enodes interface{} `json:"Enodes"`
}

type GroupIdAndEnodes struct {
	Gid    string   `json:"Gid"`
	Enodes []string `json:"Enodes"`
}
