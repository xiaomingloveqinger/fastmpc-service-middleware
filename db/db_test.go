package db

import (
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"strings"
	"testing"
)

func TestDialect_GetStructValue(t *testing.T) {
	type test struct {
		Id         int
		Name       string
		Age        int
		Flag       bool
		What       string
		Floatvalue float64
	}
	ret, err := Conn.GetStructValue("select * from test", test{})
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, v := range ret {
		c := v.(*test)
		log.Info("data", "id", c.Id, "name", c.Name, "age", c.Age, "Flag", c.Flag, "What", c.What, "Floatvalue", c.Floatvalue)
	}

	ret, err = Conn.GetStructValue("select * from test where id > ? ", test{}, 1)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, v := range ret {
		c := v.(*test)
		log.Info("data", "id", c.Id, "name", c.Name, "age", c.Age, "Flag", c.Flag, "What", c.What)
	}
}

func TestDialect_GetStringValue(t *testing.T) {
	ret, err := Conn.GetStringValue("select name from test where id = ?", 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info("data", "name", ret)

	ret, err = Conn.GetStringValue("select age from test where id = ?", 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info("data", "age string type", 31)
}

func TestDialect_GetIntValue(t *testing.T) {
	ret, err := Conn.GetIntValue("select age from test where id = ?", 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info("data", "age", ret)

	ret, err = Conn.GetIntValue("select name from test where id = ?", 2)
	if err != nil {
		if strings.Contains(err.Error(), "strconv.Atoi: parsing") {
			log.Info("ok")
		}
	}
}

func TestCommitOneRow(t *testing.T) {
	affected, err := Conn.CommitOneRow("insert into test(name,age,flag) values(?,?,?)", []interface{}{"test111", 111, 1}...)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info("insertOneRow", "affected", affected)
}

func TestBatchBatch(t *testing.T) {
	tx, err := Conn.Begin()
	if err != nil {
		t.Fatal(err.Error())
	}
	affected, err := BatchExecute("insert into test(name,age,flag) values(?,?,?)", tx, []interface{}{"test111", 111, 1}...)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info("BatchExecute", "affected", affected)
	err = Conn.Commit(tx)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGetFloatValue(t *testing.T) {
	f1, err := Conn.GetFloatValue("select floatvalue from test where id = ?", 1)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Info("getfloatvalue", "value", f1)

	f1, err = Conn.GetFloatValue("select name from test where id = ?", 1)
	if err != nil {
		if strings.Contains(err.Error(), "strconv.ParseFloat: parsing") {
			log.Info("OK")
		} else {
			t.Fatal("unexpected error")
		}
	}
}

func init() {
	common.Conf.DbConfig.DbDriverSource = "root:12345678@tcp(127.0.0.1:3306)/smw"
	common.Conf.DbConfig.DbDriverName = "mysql"
	Init()
}
