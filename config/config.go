package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type Config struct {
	Server   string `ini:"server"`
	Username string `ini:"username"`
	Passwd   string `ini:"passwd"`
	Team     string `ini:"team"`
	Chname   string `ini:"chname"`
}

func (c Config) String() string {
	server := fmt.Sprintf("[Server:%v]/[Username:%v]/[Passwd:%v]/[Team:%v]/[Chname:%v]", c.Server, c.Username, c.Passwd, c.Team, c.Chname)
	return server
}

//Read Server's Config Value from "path"
func ReadConfig(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path)
	if err != nil {
		fmt.Println("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config)
	if err != nil {
		fmt.Println("mapto config file fail!")
		return config, err
	}
	return config, nil
}
