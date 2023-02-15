package service

import "github.com/anyswap/fastmpc-service-middleware/db"

func GetTestData() (string, error) {
	return db.Conn.GetStringValue("select name from test where id = ?", 1)
}
