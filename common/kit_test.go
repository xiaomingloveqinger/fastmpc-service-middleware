package common

import (
	"errors"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"testing"
)

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
