package common

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"testing"
)

func TestConvertArrStrToStr(t *testing.T) {
	str := []string{"ab", "cd", "ef"}
	println(ConvertArrStrToStr(str))
}

func TestRecoverAddress(t *testing.T) {
	addr, err := RecoverAddress("33f0a90258905c79ba4fa9dbfa3dd166f4662ebd7f4d5dd7279926419741ea43e90ec9d27462565ea99a6c49272fe5e9704265e05083a38d60373fc5907338b0", "0x44d5b27d42f47c86233e38ff266d9e91759e1a1f87378044767b074558d0ff5c777c6d2e7d7d0f60fc7826d4fc96cbd46d8856e466b825e2c7aa54cdf927963201")
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info(addr)
}

func TestCheckThreshold(t *testing.T) {
	s := "1#2/21"
	_, _, err := CheckThreshold(s)
	if err != nil {
		log.Info(err.Error())
	}

	s = "22/21"
	_, _, err = CheckThreshold(s)
	if err != nil {
		log.Info(err.Error())
	}

	s = "22@21"
	_, _, err = CheckThreshold(s)
	if err != nil {
		log.Info(err.Error())
	}

	s = "0/21"
	_, _, err = CheckThreshold(s)
	if err != nil {
		log.Info(err.Error())
	}

	s = "1/21"
	p1, p2, err := CheckThreshold(s)
	if err != nil {
		t.Fatal(errors.New("Unexpected error"))
	}
	log.Info("test data", "p1", p1, "p2", p2)

}

func TestCheckUserAccountsAndIpPortAddr(t *testing.T) {
	s := "0x89b36c41175bc9f2341b24c8083633c89b144023|10.40.210.253"
	var param []string
	param = append(param, s)
	_, _, err := CheckUserAccountsAndIpPortAddr(param)
	if err != nil {
		log.Info(err.Error())
	}
	s1 := "0x89b36c41175bc9f2341b24c8083633c89b144023|10.40.210.253:1022"
	s2 := "0x89b36c41175bc9f2341b24c8083633c89b1440232"
	param = []string{}
	param = append(param, s1)
	param = append(param, s2)
	_, _, err = CheckUserAccountsAndIpPortAddr(param)
	if err != nil {
		log.Info(err.Error())
	}
	s3 := "0x89b36c41175bc9f2341b24c8083633c89b144023"
	param = []string{}
	param = append(param, s1)
	param = append(param, s3)
	acs, ports, err := CheckUserAccountsAndIpPortAddr(param)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, v := range acs {
		log.Info(v)
	}
	for _, v := range ports {
		log.Info(v)
	}
}

func TestGetRandomIndex(t *testing.T) {
	println(GetRandomIndex(100))
	println(GetRandomIndex(100))
	println(GetRandomIndex(100))
	println(GetRandomIndex(100))
}

func TestRawData(t *testing.T) {
	type Msg struct {
		Threshold                 string
		UserAccountsAndIpPortAddr []string
	}

	msg := &Msg{
		Threshold:                 "2/2",
		UserAccountsAndIpPortAddr: []string{"0x24afFAe9C683b7615D4130300288E348E4b5D091|127.0.0.1:3794", "0xe3fF6ebf6F0B76dC7B029d3CCB8c9eaF3F979Cff"},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info(hex.EncodeToString(data))
}

func TestPublicKeyBytesToAddress(t *testing.T) {
	pub := "04bf1b0b1a551dbfcc94dd22d605ae5c8bc591f3a6543c2878dbc6e0428ded17d68f223fb396f2773bba712b4bfbedf71ccff6c73ba8d3b9e93531a237277380d9"
	buf, err := hex.DecodeString(pub)
	if err != nil {
		t.Fatal(err)
	}
	println(PublicKeyBytesToAddress(buf).String())
}

func TestStripEnode(t *testing.T) {
	s := "enode://748ba7475b0da18887480871eb6a41c0b207c2056bf9e0cbe2d25677fef9849e3ec82d038e3d820ba9586abd1a1327555c63c34b71d9b8bccd7a1e3bedeca47b@127.0.0.1:30823"
	println(StripEnode(s))
}
