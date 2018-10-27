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
	Orgjpg   string `ini:"orgjpg"`
	Newjpg   string `ini:"newjpg"`
	Distance int    `ini:"distance"`
}

func (c Config) String() string {
	server := fmt.Sprintf("[Server:%v]/[Username:%v]/[Passwd:%v]/[Team:%v]/[Chname:%v]/[Org:%v]/[New:%v]/[Distence:%v]", c.Server, c.Username, c.Passwd, c.Team, c.Chname, c.Orgjpg, c.Newjpg, c.Distance)
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
