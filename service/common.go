package service

type SignHistory struct {
	User_account    string
	Group_id        string
	Key_id          string
	Key_type        string
	Mode            string
	Msg_context     []string
	Msg_hash        []string
	Public_key      string
	Mpc_address     string
	Threshold       string
	Txid            string
	Status          string
	Reply_Status    string
	Reply_Timestamp string
}

type SignCurNodeInfo struct {
	User_account string
	Group_id     string
	Key_id       string
	Key_type     string
	Mode         string
	Msg_context  []string
	Msg_hash     []string
	Nonce        string
	Public_key   string
	Mpc_address  string
	Threshold    string
	TimeStamp    string
	Status       int
}

type account_enode struct {
	Account string `json:"Account"`
	Enode   string `json:"Enode"`
	Ip_port string `json:"Ip_port"`
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

type RespAddr struct {
	Status          string `json:"Status"`
	User_account    string `json:"User_account"`
	Key_id          string `json:"Key_id"`
	Public_key      string `json:"Public_key"`
	Mpc_address     string `json:"Mpc_address"`
	Initializer     string `json:"Initializer"`
	Reply_status    string `json:"Reply_status"`
	Reply_timestamp string `json:"Reply_timestamp"`
	Reply_enode     string `json:"Reply_enode"`
	Gid             string `json:"Gid"`
	Threshold       string `json:"Threshold"`
}

type TxDataSign struct {
	TxType        string
	Account       string
	Nonce         string
	PubKey        string
	InputCode     string
	MsgHash       []string
	MsgContext    []string
	Keytype       string
	GroupID       string
	ThresHold     string
	Mode          string
	AcceptTimeOut string
	TimeStamp     string
	FixedApprover []string
	Comment       string
	ChainType     int
}

type Account struct {
	User_account string
	Enode        string
	Ip_port      string
}

type AcceptSignData struct {
	TxType     string   `json:"TxType"`
	Account    string   `json:"Account"`
	Nonce      string   `json:"Nonce"`
	Key        string   `json:"Key"`
	Accept     string   `json:"Accept"`
	MsgHash    []string `json:"MsgHash"`
	MsgContext []string `json:"MsgContext"`
	TimeStamp  string   `json:"TimeStamp"`
	ChainType  int      `json:"ChainType"`
}
