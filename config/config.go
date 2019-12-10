package config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	Token string
	BotPrefix string
	MysqlUser string
	MysqlPass string
	MysqlHost string
	MysqlDatabase string
	MysqlPort string

	config *configStruct
)

type configStruct struct {
	Token string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
	mysqlUser string `json:"mysqlUser"`
	mysqlPass string `json:"mysqlPass"`
	mysqlHost string `json:"mysqlHost"`
	mysqlDatabase string `json:"mysqlDB"`
	mysqlPort string `json:"mysqlPort"`
}

func ReadConfig() error {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err.Error())
	}

	Token = config.Token
	BotPrefix = config.BotPrefix
	MysqlUser = config.mysqlUser
	MysqlPass = config.mysqlPass
	MysqlDatabase = config.mysqlDatabase
	MysqlPort = config.mysqlPort
	MysqlHost = config.mysqlHost

	return nil
}