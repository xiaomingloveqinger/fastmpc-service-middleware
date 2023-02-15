package common

import (
	"encoding/json"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"os"
)

type config struct {
	Db DB
}

type DB struct {
	DbDriverName   string
	DbDriverSource string
}

var (
	Conf       = new(config)
	Configfile string
)

func Init() {
	name := "config.json"
	if Configfile != "" {
		name = Configfile
	}
	if _, err := os.Stat(name); err != nil {
		log.Error("config file is not exist", "name", name)
		os.Exit(0)
	}
	file, _ := os.Open(name)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Conf)
	if err != nil {
		log.Error("Error init Config :", "msg", err.Error())
		os.Exit(0)
	}
}
