package cache

import (
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"strconv"
	"testing"
)

func TestRCache_SetGetValue(t *testing.T) {
	key := "first"
	err := Cache.SetValue(key, 100, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, err := Cache.GetValue(key)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info(ret)
}

func TestRCache_DeleteValue(t *testing.T) {
	key := "first"
	err := Cache.SetValue(key, 100, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = Cache.DeleteValue(key)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = Cache.GetValue(key)
	if err != nil {
		log.Info(err.Error())
	}
	err = Cache.DeleteValue("sec")
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestRCache_DeleteValueByPrefix(t *testing.T) {
	key := "first"
	for i := 0; i < 10; i++ {
		err := Cache.SetValue(key+strconv.Itoa(i), 100, 0)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	err := Cache.DeleteValueByPrefix(key)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestRCache_SetJsonGetJsonValue(t *testing.T) {
	type people struct {
		Name string
		Age  int
	}
	p := people{
		Name: "clark",
		Age:  10,
	}
	err := Cache.SetJsonValue("people", &p, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
	rp := &people{}
	err = Cache.GetJsonValue("people", rp)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info("people json data", "name", rp.Name, "age", rp.Age)
}

func init() {
	common.Conf.RedisConfig.Addr = "localhost:6379"
	common.Conf.RedisConfig.Password = ""
	common.Conf.RedisConfig.DB = 0
	common.Conf.RedisConfig.PoolSize = 100
	Init()
}
