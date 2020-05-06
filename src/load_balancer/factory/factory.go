package factory

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var LBConfig Config

func checkErr(err error) {
	if err != nil {
		err = fmt.Errorf("[Configuration] %s", err.Error())
		log.Panic(err.Error())
	}
}

func InitConfigFactory(f string) {
	content, err := ioutil.ReadFile(f)
	checkErr(err)

	LBConfig = Config{}

	err = yaml.Unmarshal([]byte(content), &LBConfig)
	checkErr(err)
}
