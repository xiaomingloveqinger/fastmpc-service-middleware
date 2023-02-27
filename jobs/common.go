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

type Data struct {
	Ip_port string
	Key_id  string
	Uuid    string
}
