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
	Enode     string
	Approver  string
	Status    string
	TimeStamp string
	Initiator string
}

type Data struct {
	Ip_port string
	Key_id  string
	Uuid    string
}
