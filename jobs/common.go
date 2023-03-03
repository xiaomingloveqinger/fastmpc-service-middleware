package jobs

type ReqAddrStatus struct {
	Status    string     `json:"Status"`
	PubKey    string     `json:"PubKey"`
	Tip       string     `json:"Tip"`
	Error     string     `json:"Error"`
	AllReply  []AllReply `json:"AllReply"`
	TimeStamp string     `json:"TimeStamp"`
}

type AllReply struct {
	Enode     string `json:"Enode"`
	Approver  string `json:"Approver"`
	Status    string `json:"Status"`
	TimeStamp string `json:"timeStamp"`
	Initiator string `json:"Initiator"`
}

type SignStatus struct {
	KeyID      string     `json:"KeyID"`
	From       string     `json:"From"`
	GroupID    string     `json:"GroupID"`
	ThresHold  string     `json:"ThresHold"`
	MsgHash    []string   `json:"MsgHash"`
	MsgContext []string   `json:"MsgContext"`
	Status     string     `json:"Status"`
	Rsv        []string   `json:"Rsv"`
	Tip        string     `json:"Tip"`
	Error      string     `json:"Error"`
	Timestamp  string     `json:"Timestamp"`
	Initiator  string     `json:"Initiator"`
	PubKey     string     `json:"PubKey"`
	Keytype    string     `json:"Keytype"`
	Mode       string     `json:"Mode"`
	AllReply   []AllReply `json:"AllReply"`
	TimeStamp  string     `json:"TimeStamp"`
}

type Data struct {
	Ip_port string
	Key_id  string
	Uuid    string
}

type SignKids struct {
	Key_id string
	Pubkey string
}

type IpData struct {
	Ip_port      string
	User_account string
	Enode        string
}
