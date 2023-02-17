package common

import (
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
