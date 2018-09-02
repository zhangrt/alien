package alien

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Port string
	Conn string
	MaxIdleConns int
	MaxOpenConns int
}

func (this *Config) init() *Config  {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err =  yaml.Unmarshal(file, this)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return this
}


