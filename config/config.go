package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DB       *DBConfig
	DataPath string
}

type DBConfig struct {
	Dialect string
	//Username string
	//Password string
	//Name     string
	//Charset  string
	//Port int
	//Host string
}

func GetConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("cant open file %v", err.Error())
		os.Exit(100)
	}
	decoder := json.NewDecoder(file)

	c := Config{}
	er := decoder.Decode(&c)
	if er != nil {
		fmt.Printf("cant open file %v", err.Error())
		os.Exit(100)
	}
	return &c
}
