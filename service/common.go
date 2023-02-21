package service

type response struct {
	Status string      `json:"Status"`
	Tip    string      `json:"Tip"`
	Error  string      `json:"Error"`
	Data   interface{} `json:"Data"`
}

type account_enode struct {
	Account string
	Enode   string
	Ip_port string
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
