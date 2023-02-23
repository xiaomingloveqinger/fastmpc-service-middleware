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

type Group struct {
	Gid  string `json:"Gid"`
	Sigs string `json:"Sigs"`
	Uuid string `json:"Uuid"`
}

type TxDataReqAddr struct {
	TxType        string
	Account       string
	Nonce         string
	Keytype       string
	GroupID       string
	ThresHold     string
	Mode          string
	FixedApprover []string
	AcceptTimeOut string
	TimeStamp     string
	Sigs          string
	Comment       string
	Uuid          string
}
